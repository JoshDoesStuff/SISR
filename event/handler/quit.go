package handler

import (
	"context"

	"github.com/Alia5/SISR/sdl"
)

func Quit(rp *RegisterParams) Operation[sdl.Event] {
	return Operation[sdl.Event]{
		Event: sdl.EventTypeQuit,
		Handler: HandleFunc(
			func(_ context.Context, _ sdl.Event) error {
				rp.QuitFn()
				return nil
			},
		),
	}
}
