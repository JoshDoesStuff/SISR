package input

import (
	"context"
	"log/slog"
	"slices"
	"sync"

	"github.com/Alia5/SISR/input/viiperdevice"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/VIIPER/apiclient"
	"github.com/Alia5/VIIPER/apitypes"
)

type ViiperBridge interface {
	CreateDevice(ctx context.Context, gamepadID sdl.GamepadID, deviceType string) (chan *viiperdevice.Device, chan error)
	IsCreateDeviceScheduled(gamepadID sdl.GamepadID) bool
	Ping(ctx context.Context) (*apitypes.PingResponse, error)
}

// const minSupportedVIIPERVersion = "v0.6.1"
const defaultDeviceType = "xbox360"

func NewViiperBridge(ctx context.Context, dl DeviceStore) ViiperBridge {
	v := &viiperBridge{
		client:     apiclient.New("localhost:3242"),
		scheduled:  make(map[sdl.GamepadID]*createScheduleWhatever),
		deviceList: dl,
	}
	return v
}

type viiperBridge struct {
	client           *apiclient.Client
	busID            uint32
	viiperServerInfo *apitypes.PingResponse
	deviceList       DeviceStore

	scheduled map[sdl.GamepadID]*createScheduleWhatever

	mtx sync.Mutex
}

type createScheduleWhatever struct {
	deviceType   string
	createdChans []chan *viiperdevice.Device
	errorChans   []chan error
}

func (v *viiperBridge) CreateDevice(ctx context.Context, gamepadID sdl.GamepadID, deviceType string) (chan *viiperdevice.Device, chan error) {
	v.mtx.Lock()
	defer v.mtx.Unlock()

	deviceChan := make(chan *viiperdevice.Device, 1)
	errorChan := make(chan error, 1)
	if meh, ok := v.scheduled[gamepadID]; ok {
		meh.deviceType = deviceType
		meh.createdChans = append(meh.createdChans, deviceChan)
		meh.errorChans = append(meh.errorChans, errorChan)
		return deviceChan, errorChan
	} else {
		v.scheduled[gamepadID] = &createScheduleWhatever{
			deviceType:   deviceType,
			createdChans: []chan *viiperdevice.Device{deviceChan},
			errorChans:   []chan error{errorChan},
		}
	}
	go func() {
		defer func() {
			v.mtx.Lock()
			defer v.mtx.Unlock()
			delete(v.scheduled, gamepadID)
		}()

		// TODO: wait for ready

		busID, err := v.ensureBus(ctx)
		if err != nil {
			v.mtx.Lock()
			defer v.mtx.Unlock()
			for _, ch := range v.scheduled[gamepadID].errorChans {
				ch <- err
			}
			return
		}

		v.mtx.Lock()
		deviceType := v.scheduled[gamepadID].deviceType
		v.mtx.Unlock()
		if deviceType == "" {
			deviceType = defaultDeviceType
		}

		s, r, err := v.client.AddDeviceAndConnect(ctx, busID, deviceType, nil)
		if err != nil {
			v.mtx.Lock()
			defer v.mtx.Unlock()
			for _, ch := range v.scheduled[gamepadID].errorChans {
				ch <- err
			}
			return
		}
		closeFunc := func() error {
			_, err := v.client.DeviceRemoveCtx(context.Background(), r.BusID, r.DevId)
			return err
		}
		vd := viiperdevice.New(s, r, closeFunc)
		v.mtx.Lock()
		defer v.mtx.Unlock()
		for _, ch := range v.scheduled[gamepadID].createdChans {
			ch <- vd
		}
	}()
	return deviceChan, errorChan
}

func (v *viiperBridge) IsCreateDeviceScheduled(gamepadID sdl.GamepadID) bool {
	v.mtx.Lock()
	defer v.mtx.Unlock()

	_, ok := v.scheduled[gamepadID]
	return ok
}

func (v *viiperBridge) Ping(ctx context.Context) (*apitypes.PingResponse, error) {
	resp, err := v.client.PingCtx(ctx)
	if err != nil {
		return nil, err
	}
	v.mtx.Lock()
	defer v.mtx.Unlock()
	v.viiperServerInfo = resp
	return resp, nil
}

func (v *viiperBridge) Ready() bool {
	v.mtx.Lock()
	defer v.mtx.Unlock()

	return v.viiperServerInfo != nil
}

func (v *viiperBridge) ensureBus(ctx context.Context) (busID uint32, err error) {
	v.mtx.Lock()
	defer v.mtx.Unlock()

	if v.busID == 0 {
		slog.Debug("No previews VIIPER bus used, creating new bus")
		bus, err := v.client.BusCreateCtx(ctx, 0)
		if err != nil {
			return 0, err
		}
		v.busID = bus.BusID
		slog.Info("Created VIIPER bus", "busID", v.busID)
		return v.busID, nil
	}

	busResp, err := v.client.BusListCtx(ctx)
	if err != nil {
		return 0, err
	}
	if !slices.Contains(busResp.Buses, v.busID) {
		slog.Warn("Created VIIPER bus not found; Recreating")
		bus, err := v.client.BusCreateCtx(ctx, v.busID)
		if err != nil {
			return 0, err
		}
		v.busID = bus.BusID
		slog.Info("Re-Created VIIPER bus", "busID", v.busID)
	}
	return v.busID, nil
}
