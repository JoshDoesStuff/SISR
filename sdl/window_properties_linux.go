//go:build linux

package sdl

/*
#include <stdlib.h>
#include <SDL3/SDL_video.h>
#include <SDL3/SDL_properties.h>
*/
import "C"

import "unsafe"

// GetPointerProperty gets a platform-specific native window pointer property.
func (w *Window) GetPointerProperty(property WindowPointerProperty) uintptr {
	props := C.SDL_GetWindowProperties(w.cWindow)
	if props == 0 {
		return 0
	}

	var key string
	switch property {
	case WindowPointerPropertyWaylandDisplay:
		key = C.SDL_PROP_WINDOW_WAYLAND_DISPLAY_POINTER
	case WindowPointerPropertyWaylandSurface:
		key = C.SDL_PROP_WINDOW_WAYLAND_SURFACE_POINTER
	case WindowPointerPropertyX11Display:
		key = C.SDL_PROP_WINDOW_X11_DISPLAY_POINTER
	default:
		return 0
	}

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	p := C.SDL_GetPointerProperty(props, cKey, nil)
	return uintptr(p)
}

// GetX11WindowNumber gets the native X11 window ID for this SDL window.
// Returns 0 when unavailable (non-X11 backend or missing property).
func (w *Window) GetX11WindowNumber() uintptr {
	props := C.SDL_GetWindowProperties(w.cWindow)
	if props == 0 {
		return 0
	}

	key := C.CString(C.SDL_PROP_WINDOW_X11_WINDOW_NUMBER)
	defer C.free(unsafe.Pointer(key))

	n := C.SDL_GetNumberProperty(props, key, 0)
	if n <= 0 {
		return 0
	}

	return uintptr(n)
}
