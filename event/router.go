package event

import (
	"context"
	"log/slog"
	"runtime"

	"github.com/Alia5/SISR/event/handler"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/sdl/extras"
	"github.com/Alia5/SISR/webview"
)

type Router interface {
	RouteEvent(ctx context.Context, sdlEvent sdl.Event)
}

type router struct {
	window   *sdl.Window
	renderer sdl.Renderer
	webview  webview.WebView
	quitFn   func()

	handlers map[sdl.EventType][]func(context.Context, sdl.Event) error
}

func NewRouter(window *sdl.Window, renderer sdl.Renderer, wv webview.WebView, quitFn func()) Router {
	r := &router{
		window:   window,
		renderer: renderer,
		webview:  wv,
		quitFn:   quitFn,
		handlers: make(map[sdl.EventType][]func(context.Context, sdl.Event) error),
	}
	r.RegisterHandlers(&handler.RegisterParams{
		Window:  window,
		WebView: wv,
		QuitFn:  quitFn,
	})
	return r
}

func (r *router) RouteEvent(appCtx context.Context, ev sdl.Event) {
	if runtime.GOOS == "linux" {
		err := extras.HandleCursorHitTestWindowEvent(r.window, ev)
		if err != nil {
			slog.Error("Failed to handle cursor hit test window event", "error", err)
		}
	}
	ctx := context.WithValue(appCtx, "eventTime", ev.Base().Timestamp)
	handlers := r.handlers[ev.Base().Type]
	if len(handlers) > 0 {
		for _, handler := range handlers {
			err := handler(ctx, ev)
			if err != nil {
				slog.Error("Failed to handle event", "error", err)
			}
		}
	}
}
