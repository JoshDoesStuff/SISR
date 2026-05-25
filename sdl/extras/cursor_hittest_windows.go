//go:build windows

package extras

import (
	"log/slog"

	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/windows"
)

// GetWindowHitTest returns whether cursor hit-testing is currently enabled.
func GetWindowHitTest(window *sdl.Window) bool {
	hwnd := window.GetPointerProperty(sdl.WindowPointerPropertyWin32HWND)
	if hwnd == 0 {
		return true
	}

	hasTransparent, err := windows.HasWindowExStyleBits(hwnd, windows.WSExTransparent)
	if err != nil {
		return true
	}

	return !hasTransparent
}

// SetCursorHitTest controls whether the SDL window receives mouse hit-testing.
func SetCursorHitTest(window *sdl.Window, hittest bool) error {
	hwnd := window.GetPointerProperty(sdl.WindowPointerPropertyWin32HWND)
	if hwnd == 0 {
		return nil
	}
	hasLayered, err := windows.HasWindowExStyleBits(hwnd, windows.WSExLayered)
	if err != nil {
		return err
	}
	hasTransparent, err := windows.HasWindowExStyleBits(hwnd, windows.WSExTransparent)
	if err != nil {
		return err
	}
	if hittest {
		if hasTransparent {
			err = windows.UpdateWindowExStyleBits(hwnd, 0, windows.WSExTransparent)
			if err != nil {
				return err
			}
			err = windows.UpdateChildWindowsExStyleBits(hwnd, 0, windows.WSExTransparent)
			if err != nil {
				slog.Error("Failed to clear transparent style on child windows", "error", err)
			}
		}
		return nil
	}
	setBits := uintptr(windows.WSExTransparent)
	if !hasLayered {
		setBits |= windows.WSExLayered
	}
	if !hasTransparent || !hasLayered {
		err = windows.UpdateWindowExStyleBits(hwnd, setBits, 0)
		if err != nil {
			return err
		}
		err = windows.UpdateChildWindowsExStyleBits(hwnd, windows.WSExTransparent, 0)
		if err != nil {
			slog.Error("Failed to set transparent style on child windows", "error", err)
		}
	}
	return nil
}

// HandleCursorHitTestWindowEvent non linux stub
func HandleCursorHitTestWindowEvent(window *sdl.Window, event sdl.Event) error {
	_ = window
	_ = event
	return nil
}
