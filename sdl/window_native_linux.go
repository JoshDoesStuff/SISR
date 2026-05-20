//go:build linux

package sdl

import (
	"errors"
	"unsafe"
)

func (w *Window) GetNativeHandle() (unsafe.Pointer, error) {
	if w == nil || w.cWindow == nil {
		return nil, errors.New("window is nil")
	}

	x11Display := w.GetPointerProperty(WindowPointerPropertyX11Display)
	if x11Display != 0 {
		x11Window := w.GetX11WindowNumber()
		if x11Window != 0 {
			return unsafe.Pointer(x11Window), nil
		}
	}

	waylandSurface := w.GetPointerProperty(WindowPointerPropertyWaylandSurface)
	if waylandSurface != 0 {
		return unsafe.Pointer(waylandSurface), nil
	}

	return nil, errors.New("native window handle is unavailable")
}
