package handler

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/sdl"
)

func GamepadAdded(e *Env) Operation[*sdl.GamepadDeviceEvent] {
	return Operation[*sdl.GamepadDeviceEvent]{
		Event: sdl.EventTypeGamepadAdded,
		Handler: HandleFunc(func(ctx context.Context, ev *sdl.GamepadDeviceEvent) error {
			gpID := sdl.GamepadID(ev.Which)
			dev, err := e.DeviceStore.OpenGamePad(gpID)
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
			if err != nil {
				slog.Error("Failed to open gamepad", "error", err)
				return err
			}

			dev.Lock()
			defer dev.Unlock()

			if dev.SteamVirtualGamepad != nil {
				createViiperDevice(ctx, e.ViiperBridge, gpID, dev)
			}

			return nil
		}),
	}
}
