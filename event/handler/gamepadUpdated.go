package handler

import (
	"context"
	"log/slog"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/sdl"
)

func GamepadUpdated(c *cmd.SISRContext) Operation[*sdl.GamepadDeviceEvent] {
	return Operation[*sdl.GamepadDeviceEvent]{
		Event:   sdl.EventTypeGamepadUpdateComplete,
		Handler: HandleFunc(gpUpdate(c)),
	}
}

func gpUpdate(c *cmd.SISRContext) func(ctx context.Context, ev *sdl.GamepadDeviceEvent) error {
	return func(ctx context.Context, ev *sdl.GamepadDeviceEvent) error {
		gpID := sdl.GamepadID(ev.Which)
		dev, ok := c.DeviceStore.DeviceForID(gpID)
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
			slog.Debug("No VIIPER device for gamepad found, scheduling create", "id", gpID)
			createViiperDevice(ctx, c, gpID, dev)
			return nil
		}
		if dev.ViiperDevice.IsClosed() {
			slog.Info("VIIPER device is closed, cleaning up...", "id", gpID)
			// clear, exit, is recreated on next event
			dev.ViiperDevice = nil
			return nil
		}
		dev.ViiperDevice.Update(dev.SteamVirtualGamepad)

		return nil
	}
}
