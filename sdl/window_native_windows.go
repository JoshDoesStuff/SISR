//go:build windows

package sdl

import (
	"errors"
	"unsafe"
)

func (w *Window) GetNativeHandle() (unsafe.Pointer, error) {
	if w == nil || w.cWindow == nil {
		return nil, errors.New("window is nil")
	}

	hwnd := w.GetPointerProperty(WindowPointerPropertyWin32HWND)
	if hwnd == 0 {
		return nil, errors.New("native window handle is unavailable")
	}

	return unsafe.Pointer(hwnd), nil //nolint:govet
}
