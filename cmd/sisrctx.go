package cmd

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/Alia5/SISR/config"
	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/input/steaminputbindings"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/update"
	"github.com/Alia5/SISR/webview"
)

var ErrDispatcherClosed = errors.New("window dispatcher closed before returning a result")
var ErrUnexpectedType = errors.New("window dispatcher returned unexpected result type")

type SISRContext struct {
	WindowDispatcher WindowDispatcher[any]
	DeviceStore      input.DeviceStore
	ViiperBridge     input.ViiperBridge
	BindingEnforcer  steaminputbindings.Enforcer
	QuitFn           context.CancelFunc
	UpdateChecker    update.Checker
	Config           *SessionConfig
}

type SessionConfig struct {
	*config.AutoUpdate
	*config.RunMisc
	*config.ControllerEmulation
	*config.KeyboardMouseEmulation
	*config.Viiper
	*config.Window
	*config.Steam

	mtx sync.Mutex
}

type WindowDispatcher[T any] interface {
	Schedule(f func(w *sdl.Window, wv webview.WebView) T) <-chan T
	Dispatch(w *sdl.Window, wv webview.WebView)
}

func ScheduleWindowDispatch[T any](
	ctx context.Context,
	dispatcher WindowDispatcher[any],
	fn func(w *sdl.Window, wv webview.WebView) T,
) (T, error) {
	var zero T

	resCh := dispatcher.Schedule(func(w *sdl.Window, wv webview.WebView) any {
		return fn(w, wv)
	})

	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	case resultAny, ok := <-resCh:
		if !ok {
			return zero, ErrDispatcherClosed
		}

		if resultAny == nil {
			return zero, nil
		}

		result, ok := resultAny.(T)
		if !ok {
			return zero, fmt.Errorf("%w: expected %T, got %T", ErrUnexpectedType, zero, resultAny)
		}

		return result, nil
	}
}

type scheduledWindowFunc[T any] struct {
	fn  func(w *sdl.Window, wv webview.WebView) T
	res chan T
}

type windowDispatcher[T any] struct {
	scheduled []scheduledWindowFunc[T]
	mtx       sync.Mutex
}

func NewWindowDispatcher[T any]() WindowDispatcher[T] {
	return &windowDispatcher[T]{
		scheduled: make([]scheduledWindowFunc[T], 0),
	}
}

func (d *windowDispatcher[T]) Schedule(f func(w *sdl.Window, wv webview.WebView) T) <-chan T {
	res := make(chan T, 1)
	d.mtx.Lock()
	d.scheduled = append(d.scheduled, scheduledWindowFunc[T]{fn: f, res: res})
	d.mtx.Unlock()
	return res
}

func (d *windowDispatcher[T]) Dispatch(w *sdl.Window, wv webview.WebView) {
	d.mtx.Lock()
	if len(d.scheduled) == 0 {
		d.mtx.Unlock()
		return
	}
	scheduled := d.scheduled
	d.scheduled = make([]scheduledWindowFunc[T], 0)
	d.mtx.Unlock()

	for _, f := range scheduled {
		result := f.fn(w, wv)
		f.res <- result
		close(f.res)
	}
}

func (sc *SessionConfig) Lock() {
	sc.mtx.Lock()
}

func (sc *SessionConfig) Unlock() {
	sc.mtx.Unlock()
}
