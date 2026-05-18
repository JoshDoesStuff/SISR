package sdl

/*
#cgo CFLAGS: -I${SRCDIR}/../deps/SDL/include

#include <SDL3/SDL_video.h>
*/
import "C"

import (
	"runtime/cgo"
	"unsafe"
)

//export sisrWindowHitTestGoBridge
func sisrWindowHitTestGoBridge(cWindow *C.SDL_Window, area *C.SDL_Point, callbackData unsafe.Pointer) C.SDL_HitTestResult {
	if area == nil || callbackData == nil {
		return C.SDL_HITTEST_NORMAL
	}

	h := cgo.Handle(uintptr(callbackData))
	callback, ok := h.Value().(WindowHitTestFunc)
	if !ok || callback == nil {
		return C.SDL_HITTEST_NORMAL
	}

	result := callback(&Window{cWindow: cWindow}, int(area.x), int(area.y))
	return C.SDL_HitTestResult(result)
}
