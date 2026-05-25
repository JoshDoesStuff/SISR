package cmd

import (
	"context"
	"sync"

	"github.com/Alia5/SISR/config"
	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/input/steaminputbindings"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/update"
	"github.com/Alia5/SISR/webview"
)

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
