package sisr

import (
	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/event"
	"github.com/Alia5/SISR/event/handler"
	"github.com/Alia5/SISR/event/handler/kbmforward"
)

func registerKBMForwardHandlers(r event.Router, c *cmd.SISRContext) {
	devices := &handler.KBMDevices{}

	event.RegisterHandler(r, kbmforward.KeyboardDown(c, devices))
	event.RegisterHandler(r, kbmforward.KeyboardUp(c, devices))
	event.RegisterHandler(r, kbmforward.MouseMotion(c, devices))
	event.RegisterHandler(r, kbmforward.MouseButtonDown(c, devices))
	event.RegisterHandler(r, kbmforward.MouseButtonUp(c, devices))
	event.RegisterHandler(r, kbmforward.MouseWheel(c, devices))
	event.RegisterHandler(r, kbmforward.WindowFocusLost(c, devices))
	event.RegisterHandler(r, kbmforward.WindowHidden(c, devices))
}
