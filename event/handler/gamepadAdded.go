package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/sdl"
)

func GamepadAdded(rp *RegisterParams) Operation[*sdl.GamepadDeviceEvent] {
	return Operation[*sdl.GamepadDeviceEvent]{
		Event: sdl.EventTypeGamepadAdded,
		Handler: HandleFunc(func(_ context.Context, ev *sdl.GamepadDeviceEvent) error {
			_, err := rp.DeviceHandler.OpenGamePad(sdl.GamepadID(ev.Which))
			if err == input.ErrVirtualWithoutRealGamepad || errors.Unwrap(err) == input.ErrVirtualWithoutRealGamepad {
				slog.Warn("Virtual gamepad detected without a real gamepad, Likeley XBOX / VIIPER controller; ignoring.",
					"error", err)
				return nil
			}
			if err == input.ErrVirtualAlreadyAssigned || errors.Unwrap(err) == input.ErrVirtualAlreadyAssigned {
				slog.Warn("Virtual gamepad detected without a real gamepad, Likeley XBOX / VIIPER controller; ignoring.",
					"error", err)
				return nil
			}
			return err
		}),
	}
}
