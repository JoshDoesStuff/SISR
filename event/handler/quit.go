package handler

import (
	"context"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/sdl"
)

func Quit(c *cmd.SISRContext) Operation[sdl.Event] {
	return Operation[sdl.Event]{
		Event: sdl.EventTypeQuit,
		Handler: HandleFunc(
			func(_ context.Context, _ sdl.Event) error {
				c.QuitFn()
				return nil
			},
		),
	}
}
