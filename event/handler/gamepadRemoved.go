package handler

import (
	"context"

	"github.com/Alia5/SISR/sdl"
)

func GamepadRemoved(e *Env) Operation[*sdl.GamepadDeviceEvent] {
	return Operation[*sdl.GamepadDeviceEvent]{
		Event: sdl.EventTypeGamepadRemoved,
		Handler: HandleFunc(func(_ context.Context, ev *sdl.GamepadDeviceEvent) error {
			err := e.DeviceStore.CloseGamePad(sdl.GamepadID(ev.Which))
			return err
		}),
	}
}
