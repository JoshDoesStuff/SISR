package handler

import (
	"context"
	"log/slog"

	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/sdl"
)

func createViiperDevice(ctx context.Context, env *Env, gpID sdl.GamepadID, dev *input.Device) {
	if env.ViiperBridge.IsCreateDeviceScheduled(gpID) {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	deviceChan, errChan := env.ViiperBridge.CreateDevice(ctx, gpID, "xbox360")
	go func() {
		select {
		case vd := <-deviceChan:
			dev.Lock()
			dev.SetViiperDevice(vd)
			dev.Unlock()
			slog.Info("VIIPER device created and assigned to gamepad", "gamepad_id", gpID, "viiper_device", vd.Info())
			// TODO: check settings and stuff
			err := env.BindingEnforcer.ForceOwnAppID()
			if err != nil {
				slog.Error("Failed to force SteamInput layout", "error", err)
			}
		case err := <-errChan:
			slog.Error("Failed to create VIIPER device", "error", err)
			// TODO: handle error
		case <-ctx.Done():
			slog.Error("Timed out creating VIIPER device")
		}
		cancel()
	}()
}
