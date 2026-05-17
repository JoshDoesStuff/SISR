//go:build linux

package extras

/*
#cgo linux pkg-config: wayland-client x11 xext
#include <stdint.h>
#include "cursor_hittest_linux.h"
*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"

	"github.com/Alia5/SISR/sdl"
)

var (
	cursorHitTestStateMu sync.Mutex
	cursorHitTestState   = map[*sdl.Window]bool{}

	cursorHitTestResizeCallbackMu sync.Mutex
	cursorHitTestResizeCallback   CursorHitTestResizeCallback
)

// CursorHitTestResizeCallback runs after resize-triggered hit-test reapply.
type CursorHitTestResizeCallback func(window *sdl.Window, event *sdl.WindowEvent)

// SetCursorHitTestResizeCallback sets the optional resize reapply callback.
func SetCursorHitTestResizeCallback(callback CursorHitTestResizeCallback) {
	cursorHitTestResizeCallbackMu.Lock()
	defer cursorHitTestResizeCallbackMu.Unlock()
	cursorHitTestResizeCallback = callback
}

// HandleCursorHitTestWindowEvent reapplies stored hit-test state on resize events.
func HandleCursorHitTestWindowEvent(window *sdl.Window, event sdl.Event) error {
	we, ok := event.(*sdl.WindowEvent)
	if !ok {
		return nil
	}

	if we.Type != sdl.EventTypeWindowResized && we.Type != sdl.EventTypeWindowPixelSizeChanged {
		return nil
	}

	cursorHitTestStateMu.Lock()
	hittest, hasState := cursorHitTestState[window]
	cursorHitTestStateMu.Unlock()
	if !hasState {
		return nil
	}

	if err := SetCursorHitTest(window, hittest); err != nil {
		return err
	}

	cursorHitTestResizeCallbackMu.Lock()
	callback := cursorHitTestResizeCallback
	cursorHitTestResizeCallbackMu.Unlock()
	if callback != nil {
		callback(window, we)
	}

	return nil
}

// SetCursorHitTest enables/disables window hit-testing.
func SetCursorHitTest(window *sdl.Window, hittest bool) error {
	chittest := C.int(0)
	if hittest {
		chittest = C.int(1)
	}

	waylandDisplay := window.GetPointerProperty(sdl.WindowPointerPropertyWaylandDisplay)
	waylandSurface := window.GetPointerProperty(sdl.WindowPointerPropertyWaylandSurface)
	if waylandDisplay != 0 && waylandSurface != 0 {
		if C.set_wayland_cursor_hittest(
			unsafe.Pointer(waylandDisplay),
			unsafe.Pointer(waylandSurface),
			chittest,
		) != 0 {
			return fmt.Errorf("failed to set Wayland cursor hit test")
		}
		cursorHitTestStateMu.Lock()
		cursorHitTestState[window] = hittest
		cursorHitTestStateMu.Unlock()
		return nil
	}

	x11Display := window.GetPointerProperty(sdl.WindowPointerPropertyX11Display)
	x11Window := window.GetX11WindowNumber()
	if x11Display != 0 && x11Window != 0 {
		if C.set_x11_cursor_hittest(
			unsafe.Pointer(x11Display),
			C.uintptr_t(x11Window),
			chittest,
		) != 0 {
			return fmt.Errorf("failed to set X11 cursor hit test")
		}
		cursorHitTestStateMu.Lock()
		cursorHitTestState[window] = hittest
		cursorHitTestStateMu.Unlock()
		return nil
	}

	return nil
}
