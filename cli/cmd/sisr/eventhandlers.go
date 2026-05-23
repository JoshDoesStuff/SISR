package sisr

import (
	"context"
	"runtime"

	"github.com/Alia5/SISR/event"
	"github.com/Alia5/SISR/event/handler"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/sdl/extras"
)

func registerEventHandlers(r event.Router, e *handler.Env) {
	if runtime.GOOS == "linux" {
		hittestfunc := handler.HandleFunc(
			func(_ context.Context, ev *sdl.WindowEvent) error {
				return extras.HandleCursorHitTestWindowEvent(e.Window, ev)
			},
		)
		event.RegisterHandler(r, handler.Operation[*sdl.WindowEvent]{
			Event:   sdl.EventTypeWindowPixelSizeChanged,
			Handler: hittestfunc,
		})
		event.RegisterHandler(r, handler.Operation[*sdl.WindowEvent]{
			Event:   sdl.EventTypeWindowResized,
			Handler: hittestfunc,
		})
	}
	event.RegisterHandler(r, handler.Quit(e))
	event.RegisterHandler(r, handler.WindowResize(e))
	event.RegisterHandler(r, handler.GamepadAdded(e))
	event.RegisterHandler(r, handler.GamepadRemoved(e))
	event.RegisterHandler(r, handler.GamepadUpdated(e))
}
