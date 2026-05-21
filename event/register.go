package event

import (
	"context"
	"fmt"

	"github.com/Alia5/SISR/event/handler"
	"github.com/Alia5/SISR/sdl"
)

func (r *router) RegisterHandlers(rp *handler.RegisterParams) {
	registerHandler(r, handler.HandleQuit(rp))
	registerHandler(r, handler.HandleResize(rp))
}

func registerHandler[T sdl.Event](r *router, ops ...*handler.EventOperation[T]) {
	for _, op := range ops {
		for _, event := range op.Events {
			r.handlers[event] = append(r.handlers[event], func(ctx context.Context, sdlEvent sdl.Event) error {
				typedEvent, ok := sdlEvent.(T)
				if !ok {
					return fmt.Errorf("unexpected event type %T for %v", sdlEvent, event)
				}

				return op.Handler(ctx, typedEvent)
			})
		}
	}
}
