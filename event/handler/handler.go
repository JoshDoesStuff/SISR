package handler

import (
	"context"

	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/webview"
)

type RegisterParams struct {
	Window        *sdl.Window
	WebView       webview.WebView
	DeviceHandler input.DeviceHandler
	QuitFn        context.CancelFunc
}

type Handler[T sdl.Event] interface {
	Handle(context.Context, T) error
}
type handler[T sdl.Event] struct {
	handleFunc HandlerFunc[T]
}

func (h *handler[T]) Handle(ctx context.Context, event T) error {
	return h.handleFunc(ctx, event)
}

type HandlerFunc[T sdl.Event] func(context.Context, T) error

func HandleFunc[T sdl.Event](handleEvent HandlerFunc[T]) Handler[T] {
	return &handler[T]{handleFunc: handleEvent}
}

type Operation[T sdl.Event] struct {
	Event   sdl.EventType
	Handler Handler[T]
}
