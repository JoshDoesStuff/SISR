package handler

import (
	"context"

	"github.com/Alia5/SISR/sdl"
)

func WindowResize(rp *Env) Operation[*sdl.WindowEvent] {
	return Operation[*sdl.WindowEvent]{
		Event: sdl.EventTypeWindowResized,
		Handler: HandleFunc(
			func(_ context.Context, ev *sdl.WindowEvent) error {
				rp.WebView.Resize(int(ev.Data1), int(ev.Data2))
				return nil
			},
		),
	}
}
