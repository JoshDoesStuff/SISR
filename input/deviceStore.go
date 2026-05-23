package input

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/Alia5/SISR/sdl"
)

type DeviceStore interface {
	OpenGamePad(id sdl.GamepadID) (*Device, error)
	CloseGamePad(id sdl.GamepadID) error
	DeviceForID(id sdl.GamepadID) (*Device, bool)
}

type deviceStore struct {
	devices map[sdl.GamepadID]*Device
	mtx     sync.Mutex
}

func NewDeviceStore() (DeviceStore, func(), error) {

	opener := &deviceStore{
		devices: make(map[sdl.GamepadID]*Device),
	}
	return opener, opener.quit, nil
}

func (ds *deviceStore) OpenGamePad(id sdl.GamepadID) (*Device, error) {
	ds.mtx.Lock()
	defer ds.mtx.Unlock()

	if dev, exists := ds.devices[id]; exists {
		return dev, nil
	}
	var dev *Device
	gp, err := sdl.OpenGamepad(id)
	if err != nil {
		return nil, err
	}
	steamHandle := gp.GetSteamHandle()
	serial := gp.Serial()
	path := gp.Path()
	gpType := gp.Type()
	realType := gp.RealType()

	slog.Debug("Opened Gamepad",
		"id", id,
		"name", gp.Name(),
		"steamHandle", steamHandle,
		"serial", serial,
		"path", path,
		"type", gpType.Name(),
		"realType", realType.Name(),
	)

	if steamHandle != 0 {
		if id == 0 {
			// Xbox controller, ignore for now
			gp.Close()
			return nil, fmt.Errorf(
				"%w: first encountered device id=%d (SteamHandle=%d) has no corresponding real gamepad",
				ErrVirtualWithoutRealGamepad, id, steamHandle,
			)
		}
		// We assume devices are opened in order
		// Real Gamepad first and Steam virtual gamepad immediately after with Id+1
		// I havent validated, but also haven't encountered anything else so far.
		var ok bool
		dev, ok = ds.devices[id-1]
		if !ok {
			// Xbox controller, ignore for now
			gp.Close()
			return nil, fmt.Errorf(
				"%w: device id=%d (SteamHandle=%d) has no corresponding real gamepad",
				ErrVirtualWithoutRealGamepad, id, steamHandle,
			)
		}
		dev.Lock()
		defer dev.Unlock()

		if dev.SteamVirtualGamepad != nil {
			gp.Close()
			return nil, ErrVirtualAlreadyAssigned
		}
		dev.SteamVirtualGamepad = gp
		slog.Info("Found Steam Virtual Gamepad",
			"id", id,
			"name", gp.Name(),
			"steamhandle", steamHandle,
			"paired with", dev.RealGamepad.ID(),
			"paired with name", dev.RealGamepad.Name(),
		)
	} else {
		dev = &Device{RealGamepad: gp}
		slog.Debug("Opened real gamepad", "id", id, "name", gp.Name())
	}
	ds.devices[id] = dev

	return dev, nil
}

func (ds *deviceStore) CloseGamePad(id sdl.GamepadID) error {
	ds.mtx.Lock()
	defer ds.mtx.Unlock()

	dev, ok := ds.devices[id]
	if !ok {
		return nil
	}
	dev.Lock()
	defer dev.Unlock()

	gp := dev.RealGamepad
	if gp == nil {
		gp = dev.SteamVirtualGamepad
		dev.SteamVirtualGamepad = nil
	} else {
		dev.RealGamepad = nil
	}
	slog.Debug("Closing gamepad", "id", id, "name", gp.Name(), "steamHandle", gp.GetSteamHandle())
	gp.Close()

	delete(ds.devices, id)
	return nil
}

func (ds *deviceStore) DeviceForID(id sdl.GamepadID) (*Device, bool) {
	ds.mtx.Lock()
	defer ds.mtx.Unlock()

	dev, ok := ds.devices[id]
	return dev, ok
}

func (ds *deviceStore) quit() {
	for id, dev := range ds.devices {
		if dev != nil {
			dev.Lock()
			dev.Close()
			dev.Unlock()
		}
		delete(ds.devices, id)
	}

}
