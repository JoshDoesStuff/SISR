//go:build windows

package sdl

/*
#cgo CFLAGS: -I${SRCDIR}/../deps/SDL/include
#cgo LDFLAGS: -L${SRCDIR}/../deps/SDL/build/Debug -lSDL3

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
	case WindowPointerPropertyWin32HWND:
		key = C.SDL_PROP_WINDOW_WIN32_HWND_POINTER
	case WindowPointerPropertyWin32HDC:
		key = C.SDL_PROP_WINDOW_WIN32_HDC_POINTER
	case WindowPointerPropertyWin32Instance:
		key = C.SDL_PROP_WINDOW_WIN32_INSTANCE_POINTER
	default:
		return 0
	}

	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	p := C.SDL_GetPointerProperty(props, cKey, nil)
	return uintptr(p)
}
