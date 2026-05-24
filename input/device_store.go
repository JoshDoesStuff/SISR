package input

import (
	"log/slog"
	"slices"
	"sync"

	"github.com/Alia5/SISR/sdl"
)

type DeviceStore interface {
	OpenGamePad(id sdl.GamepadID) (*Device, error)
	CloseGamePad(id sdl.GamepadID) error
	DeviceForID(id sdl.GamepadID) (*Device, bool)
	IgnoreNextDevice(num int)
	Empty() bool
}

type deviceStore struct {
	devices        map[sdl.GamepadID]*Device
	deviceIdxOrder []sdl.GamepadID
	ignoreNext     int

	mtx sync.Mutex
}

func NewDeviceStore() (DeviceStore, func(), error) {

	opener := &deviceStore{
		devices:        make(map[sdl.GamepadID]*Device),
		deviceIdxOrder: make([]sdl.GamepadID, 0),
		ignoreNext:     0,
	}
	return opener, opener.quit, nil
}

func (ds *deviceStore) OpenGamePad(id sdl.GamepadID) (*Device, error) {
	ds.mtx.Lock()
	defer ds.mtx.Unlock()
	defer func() {
		if slices.Contains(ds.deviceIdxOrder, id) {
			return
		}
		ds.deviceIdxOrder = append(ds.deviceIdxOrder, id)
	}()

	if dev, exists := ds.devices[id]; exists {
		return dev, nil
	}

	if ds.ignoreNext > 0 {
		ds.ignoreNext--
		slog.Info("Skipping gamepad due to ignore counter", "id", id, "remaining", ds.ignoreNext)
		return nil, nil
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
		if len(ds.deviceIdxOrder) == 0 {
			// xbox controllers on windows are not detected twice
			// Hid-Path + (Steam Virtual via XInput), but only via XInput
			// Thus, we assign both real and steam virtual pad from the same sdl gamepad
			slog.Warn("First opened gamepad has non-zero Steam handle; likely xbox controller")
			dev = &Device{
				RealGamepad: gp,
			}
		} else {
			var ok bool
			for _, devID := range slices.Backward(ds.deviceIdxOrder) {
				lastDevice, exists := ds.devices[devID]
				if !exists || lastDevice == nil {
					continue
				}
				if lastDevice.RealGamepad != nil && lastDevice.SteamVirtualGamepad == nil {
					dev = lastDevice
					ok = true
					break
				}
			}
			if !ok || dev == nil {
				slog.Debug("Likekely xbox controller")
				dev = &Device{
					RealGamepad: gp,
				}
			}
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

	if dev.RealGamepad != nil && dev.SteamVirtualGamepad != nil && dev.RealGamepad.ID() == id && dev.SteamVirtualGamepad.ID() == id {
		slog.Info(
			"Closing gamepad",
			"id", id,
			"name", dev.RealGamepad.Name(),
			"type", "both (real + steam virtual; likely Xbox controller)",
			"steamHandle", dev.SteamVirtualGamepad.GetSteamHandle(),
		)
		dev.RealGamepad.Close()
		dev.RealGamepad = nil
		dev.SteamVirtualGamepad = nil
	} else {
		if dev.RealGamepad != nil && dev.RealGamepad.ID() == id {
			slog.Info(
				"Closing gamepad",
				"id", id,
				"name", dev.RealGamepad.Name(),
				"type", "real",
			)
			dev.RealGamepad.Close()
			dev.RealGamepad = nil
		}
		if dev.SteamVirtualGamepad != nil && dev.SteamVirtualGamepad.ID() == id {
			slog.Info(
				"Closing gamepad",
				"id", id,
				"name", dev.SteamVirtualGamepad.Name(),
				"type", "steam_virtual",
				"steamHandle", dev.SteamVirtualGamepad.GetSteamHandle(),
			)
			dev.SteamVirtualGamepad.Close()
			dev.SteamVirtualGamepad = nil
		}
	}

	delete(ds.devices, id)
	if dev.RealGamepad == nil && dev.SteamVirtualGamepad == nil {
		slog.Info("Device has no more gamepads, cleaning up...", "id", id)
		dev.Close()
	}
	return nil
}

func (ds *deviceStore) DeviceForID(id sdl.GamepadID) (*Device, bool) {
	ds.mtx.Lock()
	defer ds.mtx.Unlock()

	dev, ok := ds.devices[id]
	return dev, ok
}

func (ds *deviceStore) IgnoreNextDevice(num int) {
	ds.mtx.Lock()
	defer ds.mtx.Unlock()

	ds.ignoreNext += num
	slog.Debug("Updated ignore-next-device counter", "added", num, "total", ds.ignoreNext)
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

func (ds *deviceStore) Empty() bool {
	ds.mtx.Lock()
	defer ds.mtx.Unlock()
	return len(ds.devices) == 0
}
