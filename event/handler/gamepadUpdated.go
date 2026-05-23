package handler

import (
	"context"

	"github.com/Alia5/SISR/sdl"
)

func GamepadUpdated(rp *RegisterParams) Operation[*sdl.GamepadDeviceEvent] {
	return Operation[*sdl.GamepadDeviceEvent]{
		Event:   sdl.EventTypeGamepadUpdateComplete,
		Handler: HandleFunc(gpUpdate(rp)),
	}
}

func gpUpdate(rp *RegisterParams) func(_ context.Context, ev *sdl.GamepadDeviceEvent) error {
	return func(_ context.Context, ev *sdl.GamepadDeviceEvent) error {
		gpID := sdl.GamepadID(ev.Which)
		dev, ok := rp.DeviceHandler.DeviceForID(gpID)
		if !ok {
			return nil
		}
		err := dev.UpdateViiperDevice(gpID)
		if err != nil {
			return err
		}

		return nil
	}
}
