package sisr

import (
	"log/slog"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/event"
	"github.com/Alia5/SISR/event/handler"
	"github.com/Alia5/SISR/event/handler/kbmforward"
	"github.com/Alia5/SISR/sdl"
)

func registerKBMForwardHandlers(r event.Router, c *cmd.SISRContext, window *sdl.Window) {
	devices := &handler.KBMDevices{}

	event.RegisterHandler(r, kbmforward.KeyboardDown(c, devices))
	event.RegisterHandler(r, kbmforward.KeyboardUp(c, devices))
	event.RegisterHandler(r, kbmforward.MouseMotion(c, devices))
	event.RegisterHandler(r, kbmforward.MouseButtonDown(c, devices))
	event.RegisterHandler(r, kbmforward.MouseButtonUp(c, devices))
	event.RegisterHandler(r, kbmforward.MouseWheel(c, devices))
	event.RegisterHandler(r, kbmforward.WindowFocusGained(c, window, devices))
	event.RegisterHandler(r, kbmforward.WindowFocusLost(c, window, devices))
	event.RegisterHandler(r, kbmforward.WindowHidden(c, window, devices))

	if err := window.SetWindowMouseGrab(true); err != nil {
		slog.Error("Failed to enable window mouse grab during KBM handler registration", "error", err)
	}
	if err := window.SetWindowRelativeMouseMode(true); err != nil {
		slog.Error("Failed to enable relative mouse mode during KBM handler registration", "error", err)
	}
}
