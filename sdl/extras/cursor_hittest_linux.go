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
	waylandCursorHitTestStateMu sync.Mutex
	waylandCursorHitTestState   = map[*sdl.Window]bool{}
)

// GetWindowHitTest returns whether cursor hit-testing is currently enabled.
//
// On X11 this is queried from the native input shape region.
// On Wayland, there is no protocol getter for wl_surface input region, so we
// retain the last successfully applied state for that window.
func GetWindowHitTest(window *sdl.Window) bool {
	waylandDisplay := window.GetPointerProperty(sdl.WindowPointerPropertyWaylandDisplay)
	waylandSurface := window.GetPointerProperty(sdl.WindowPointerPropertyWaylandSurface)
	if waylandDisplay != 0 && waylandSurface != 0 {
		waylandCursorHitTestStateMu.Lock()
		hittest, ok := waylandCursorHitTestState[window]
		waylandCursorHitTestStateMu.Unlock()
		if !ok {
			return true
		}
		return hittest
	}

	x11Display := window.GetPointerProperty(sdl.WindowPointerPropertyX11Display)
	x11Window := window.GetX11WindowNumber()
	if x11Display != 0 && x11Window != 0 {
		state := C.get_x11_cursor_hittest(
			unsafe.Pointer(x11Display),
			C.uintptr_t(x11Window),
		)
		if state < 0 {
			return true
		}
		return state != 0
	}

	return true
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

	hittest := GetWindowHitTest(window)

	if err := SetCursorHitTest(window, hittest); err != nil {
		return err
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
		waylandCursorHitTestStateMu.Lock()
		waylandCursorHitTestState[window] = hittest
		waylandCursorHitTestStateMu.Unlock()
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
		return nil
	}

	return nil
}
