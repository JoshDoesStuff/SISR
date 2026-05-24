package handler

import (
	"context"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/webview"
)

func WindowResize(c *cmd.SISRContext, webview webview.WebView) Operation[*sdl.WindowEvent] {
	return Operation[*sdl.WindowEvent]{
		Event: sdl.EventTypeWindowResized,
		Handler: HandleFunc(
			func(_ context.Context, ev *sdl.WindowEvent) error {
				webview.Resize(int(ev.Data1), int(ev.Data2))
				return nil
			},
		),
	}
}
