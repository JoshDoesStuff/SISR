//go:build windows

package windows

import (
	"fmt"
	"log/slog"
	"syscall"
	"unsafe"

	"github.com/Alia5/SISR/sdl"
)

var (
	dwmapiDLL                 = syscall.NewLazyDLL("dwmapi.dll")
	dwmSetWindowAttributeProc = dwmapiDLL.NewProc("DwmSetWindowAttribute")
)

const (
	dwmwaPassiveUpdateMode = 37
	eInvalidArgHRESULT     = 0x80070057
)

func SetDWMPassiveUpdateMode(window *sdl.Window) error {
	hwnd := window.GetPointerProperty(sdl.WindowPointerPropertyWin32HWND)
	if hwnd == 0 {
		return nil
	}

	slog.Debug("Setting DWM passive mode", "HWND", hwnd)

	passiveUpdateMode := int32(1)
	hr, _, callErr := dwmSetWindowAttributeProc.Call(
		hwnd,
		uintptr(dwmwaPassiveUpdateMode),
		uintptr(unsafe.Pointer(&passiveUpdateMode)),
		unsafe.Sizeof(passiveUpdateMode),
	)
	if hr != 0 {
		if uint32(hr) == eInvalidArgHRESULT {
			return nil
		}
		if callErr != syscall.Errno(0) {
			return fmt.Errorf("DwmSetWindowAttribute failed, HRESULT=0x%x: %w", uint32(hr), callErr)
		}
		return fmt.Errorf("DwmSetWindowAttribute failed, HRESULT=0x%x", uint32(hr))
	}

	return nil
}
