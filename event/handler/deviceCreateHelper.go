package handler

import (
	"context"
	"log/slog"
	"time"

	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/sdl"
)

func createViiperDevice(ctx context.Context, env *Env, gpID sdl.GamepadID, dev *input.Device) {
	if !env.ViiperBridge.Ready() {
		return
	}
	if env.ViiperBridge.IsCreateDeviceScheduled(gpID) {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	deviceChan, errChan := env.ViiperBridge.CreateDevice(ctx, gpID, env.Config.DefaultControllerType)
	go func() {
		select {
		case vd := <-deviceChan:
			ignoreNextCount := 1
			if vd.Info().Type != "xbox360" {
				ignoreNextCount = 2
			}
			env.DeviceStore.IgnoreNextDevice(ignoreNextCount)
			dev.Lock()
			dev.SetViiperDevice(vd)
			dev.Unlock()
			slog.Info("VIIPER device created and assigned to gamepad", "gamepad_id", gpID, "viiper_device", vd.Info())
		case err := <-errChan:
			slog.Error("Failed to create VIIPER device", "error", err)
			// best attempt... otherwise user will have handle to via UI
			// kiss
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			env.ViiperBridge.Ping(ctx) // nolint
		case <-ctx.Done():
			slog.Error("Timed out creating VIIPER device")
		}
		cancel()
	}()
}
