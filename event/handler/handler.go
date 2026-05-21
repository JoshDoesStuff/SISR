package handler

import (
	"context"

	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/webview"
)

type RegisterParams struct {
	Window  *sdl.Window
	WebView webview.WebView
	QuitFn  context.CancelFunc
}

type HandlerFunc[T sdl.Event] func(context.Context, T) error

func HandleFunc[T sdl.Event](handleEvent HandlerFunc[T]) HandlerFunc[T] {
	return handleEvent
}

type EventOperation[T sdl.Event] struct {
	Events  []sdl.EventType
	Handler HandlerFunc[T]
}
