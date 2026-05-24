package handler

import (
	"context"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/sdl"
)

func GamepadRemoved(c *cmd.SISRContext) Operation[*sdl.GamepadDeviceEvent] {
	return Operation[*sdl.GamepadDeviceEvent]{
		Event: sdl.EventTypeGamepadRemoved,
		Handler: HandleFunc(func(_ context.Context, ev *sdl.GamepadDeviceEvent) error {
			err := c.DeviceStore.CloseGamePad(sdl.GamepadID(ev.Which))
			return err
		}),
	}
}
