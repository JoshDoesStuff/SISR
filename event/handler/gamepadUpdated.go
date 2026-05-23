package handler

import (
	"context"

	"github.com/Alia5/SISR/sdl"
)

func GamepadUpdated(e *Env) Operation[*sdl.GamepadDeviceEvent] {
	return Operation[*sdl.GamepadDeviceEvent]{
		Event:   sdl.EventTypeGamepadUpdateComplete,
		Handler: HandleFunc(gpUpdate(e)),
	}
}

func gpUpdate(e *Env) func(ctx context.Context, ev *sdl.GamepadDeviceEvent) error {
	return func(ctx context.Context, ev *sdl.GamepadDeviceEvent) error {
		gpID := sdl.GamepadID(ev.Which)
		dev, ok := e.DeviceStore.DeviceForID(gpID)
		if !ok {
			return nil
		}
		dev.Lock()
		defer dev.Unlock()

		if dev.SteamVirtualGamepad == nil {
			return nil
		}
		if dev.SteamVirtualGamepad.ID() != gpID {
			return nil
		}

		if dev.ViiperDevice == nil {
			createViiperDevice(ctx, e, gpID, dev)
		}

		dev.ViiperDevice.Update(dev.SteamVirtualGamepad)

		return nil
	}
}
