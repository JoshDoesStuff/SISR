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
		err := windows.UpdateWindowExStyleBits(hwnd, 0, windows.WSExTransparent|windows.WSExLayered)
		if err != nil {
			return err
		}
		_ = windows.UpdateChildWindowsExStyleBits(hwnd, 0, windows.WSExTransparent)
		return nil
	}
	err := windows.UpdateWindowExStyleBits(hwnd, windows.WSExTransparent|windows.WSExLayered, 0)
	if err != nil {
		return err
	}
	_ = windows.UpdateChildWindowsExStyleBits(hwnd, windows.WSExTransparent, windows.WSExLayered)
	return nil
}

// HandleCursorHitTestWindowEvent non linux stub
func HandleCursorHitTestWindowEvent(window *sdl.Window, event sdl.Event) error {
	_ = window
	_ = event
	return nil
}
