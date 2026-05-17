//go:build windows

package extras

import (
	"fmt"
	"syscall"

	"github.com/Alia5/SISR/sdl"
)

var (
	user32DLL            = syscall.NewLazyDLL("user32.dll")
	getWindowLongPtrProc = user32DLL.NewProc("GetWindowLongPtrW")
	setWindowLongPtrProc = user32DLL.NewProc("SetWindowLongPtrW")
	setWindowPosProc     = user32DLL.NewProc("SetWindowPos")
)

const (
	gwlExStyle      = -20
	wsExLayered     = 0x00080000
	wsExTransparent = 0x00000020
	swpNoSize       = 0x0001
	swpNoMove       = 0x0002
	swpNoZOrder     = 0x0004
	swpNoActivate   = 0x0010
	swpFrameChanged = 0x0020
)

// SetCursorHitTest controls whether the SDL window receives mouse hit-testing.
func SetCursorHitTest(window *sdl.Window, hittest bool) error {
	hwnd := window.GetPointerProperty(sdl.WindowPointerPropertyWin32HWND)
	if hwnd == 0 {
		return nil
	}

	gwlExStyleIndex := int32(gwlExStyle)
	index := uintptr(gwlExStyleIndex)
	exStyle, _, callErr := getWindowLongPtrProc.Call(hwnd, index)
	if exStyle == 0 && callErr != syscall.Errno(0) {
		return fmt.Errorf("GetWindowLongPtrW failed: %w", callErr)
	}

	newExStyle := exStyle
	if hittest {
		newExStyle &^= wsExTransparent
	} else {
		newExStyle |= wsExLayered | wsExTransparent
	}

	if newExStyle == exStyle {
		return nil
	}

	result, _, callErr := setWindowLongPtrProc.Call(hwnd, index, newExStyle)
	if result == 0 && callErr != syscall.Errno(0) {
		return fmt.Errorf("SetWindowLongPtrW failed: %w", callErr)
	}

	const frameFlags = swpNoSize | swpNoMove | swpNoZOrder | swpNoActivate | swpFrameChanged
	result, _, callErr = setWindowPosProc.Call(hwnd, 0, 0, 0, 0, 0, uintptr(frameFlags))
	if result == 0 {
		if callErr != syscall.Errno(0) {
			return fmt.Errorf("SetWindowPos failed: %w", callErr)
		}
		return fmt.Errorf("SetWindowPos failed")
	}

	return nil
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
