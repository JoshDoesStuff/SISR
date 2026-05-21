package handler

import (
	"context"

	"github.com/Alia5/SISR/sdl"
)

func HandleResize(rp *RegisterParams) *EventOperation[*sdl.WindowEvent] {
	return &EventOperation[*sdl.WindowEvent]{
		Events: []sdl.EventType{
			sdl.EventTypeWindowResized,
		},
		Handler: HandleFunc(
			func(_ context.Context, ev *sdl.WindowEvent) error {
				rp.WebView.Resize(int(ev.Data1), int(ev.Data2))
				return nil
			},
		),
	}
}
