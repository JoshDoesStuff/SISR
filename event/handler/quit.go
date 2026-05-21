package handler

import (
	"context"

	"github.com/Alia5/SISR/sdl"
)

func HandleQuit(rp *RegisterParams) *EventOperation[sdl.Event] {
	return &EventOperation[sdl.Event]{
		Events: []sdl.EventType{
			sdl.EventTypeQuit,
		},
		Handler: HandleFunc(
			func(_ context.Context, _ sdl.Event) error {
				rp.QuitFn()
				return nil
			},
		),
	}
}
