//go:build windows

package extras

import (
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/windows"
)

// SetCursorHitTest controls whether the SDL window receives mouse hit-testing.
func SetCursorHitTest(window *sdl.Window, hittest bool) error {
	hwnd := window.GetPointerProperty(sdl.WindowPointerPropertyWin32HWND)
	if hwnd == 0 {
		return nil
	}
	if hittest {
		return windows.UpdateWindowExStyleBits(hwnd, 0, windows.WSExTransparent)
	}
	return windows.UpdateWindowExStyleBits(hwnd, windows.WSExTransparent, windows.WSExLayered)
}

// CursorHitTestResizeCallback non linux stub
type CursorHitTestResizeCallback func(window *sdl.Window, event *sdl.WindowEvent)

// SetCursorHitTestResizeCallback non linux stub
func SetCursorHitTestResizeCallback(callback CursorHitTestResizeCallback) {
	_ = callback
}

// HandleCursorHitTestWindowEvent non linux stub
func HandleCursorHitTestWindowEvent(window *sdl.Window, event sdl.Event) error {
	_ = window
	_ = event
	return nil
}
