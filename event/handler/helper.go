package handler

import (
	"context"
	"log/slog"

	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/sdl"
)

func createViiperDevice(ctx context.Context, vp input.ViiperBridge, gpID sdl.GamepadID, dev *input.Device) {
	if vp.CreateDeviceScheduled(gpID) {
		return
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	deviceChan, errChan := vp.CreateDevice(ctx, gpID, "xbox360")
	go func() {
		select {
		case vd := <-deviceChan:
			dev.Lock()
			defer dev.Unlock()
			dev.SetViiperDevice(vd)
		case err := <-errChan:
			slog.Error("Failed to create VIIPER device", "error", err)
		case <-ctx.Done():
			slog.Error("Timed out creating VIIPER device")
		}
		cancel()
	}()
}
