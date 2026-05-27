package sdl

/*
#cgo CFLAGS: -I${SRCDIR}/../deps/SDL/include

#include <stdlib.h>

#include <SDL3/SDL_video.h>
#include <SDL3/SDL_render.h>
#include <SDL3/SDL_mouse.h>

extern int sisrSetWindowRelativeMouseModeBridge(void *window, int enabled);
extern int sisrGetWindowRelativeMouseModeBridge(void *window);

extern SDL_HitTestResult sisrWindowHitTestGoBridge(SDL_Window *win, SDL_Point *area, void *userdata);

SDL_HitTestResult SDLCALL sisrWindowHitTestBridge(SDL_Window *win, const SDL_Point *area, void *userdata) {
	return sisrWindowHitTestGoBridge(win, (SDL_Point *)area, userdata);
}
*/
import "C"
import (
	"runtime/cgo"
	"unsafe"
)

// DisplayMode defines a display mode.
type DisplayMode struct {
	DisplayID              uint32
	Format                 uint32
	W                      int
	H                      int
	PixelDensity           float32
	RefreshRate            float32
	RefreshRateNumerator   int
	RefreshRateDenominator int
}

// Rect is the SDL_Rect equivalent used by window APIs.
type Rect struct {
	X int
	Y int
	W int
	H int
}

// HitTestResult represents possible return values from the SDL_HitTest callback.
type HitTestResult int

// WindowHitTestFunc is the callback type used by SetWindowHitTest.
type WindowHitTestFunc func(window *Window, x, y int) HitTestResult

// Hit-test result values.
const (
	HitTestNormal            HitTestResult = C.SDL_HITTEST_NORMAL
	HitTestDraggable         HitTestResult = C.SDL_HITTEST_DRAGGABLE
	HitTestResizeTopLeft     HitTestResult = C.SDL_HITTEST_RESIZE_TOPLEFT
	HitTestResizeTop         HitTestResult = C.SDL_HITTEST_RESIZE_TOP
	HitTestResizeTopRight    HitTestResult = C.SDL_HITTEST_RESIZE_TOPRIGHT
	HitTestResizeRight       HitTestResult = C.SDL_HITTEST_RESIZE_RIGHT
	HitTestResizeBottomRight HitTestResult = C.SDL_HITTEST_RESIZE_BOTTOMRIGHT
	HitTestResizeBottom      HitTestResult = C.SDL_HITTEST_RESIZE_BOTTOM
	HitTestResizeBottomLeft  HitTestResult = C.SDL_HITTEST_RESIZE_BOTTOMLEFT
	HitTestResizeLeft        HitTestResult = C.SDL_HITTEST_RESIZE_LEFT
)

// WindowFlags are the flags on a window.
type WindowFlags uint64

// Window flags.
const (
	WindowFlagFullscreen        WindowFlags = C.SDL_WINDOW_FULLSCREEN
	WindowFlagOpenGL            WindowFlags = C.SDL_WINDOW_OPENGL
	WindowFlagOccluded          WindowFlags = C.SDL_WINDOW_OCCLUDED
	WindowFlagHidden            WindowFlags = C.SDL_WINDOW_HIDDEN
	WindowFlagBorderless        WindowFlags = C.SDL_WINDOW_BORDERLESS
	WindowFlagResizable         WindowFlags = C.SDL_WINDOW_RESIZABLE
	WindowFlagMinimized         WindowFlags = C.SDL_WINDOW_MINIMIZED
	WindowFlagMaximized         WindowFlags = C.SDL_WINDOW_MAXIMIZED
	WindowFlagMouseGrabbed      WindowFlags = C.SDL_WINDOW_MOUSE_GRABBED
	WindowFlagInputFocus        WindowFlags = C.SDL_WINDOW_INPUT_FOCUS
	WindowFlagMouseFocus        WindowFlags = C.SDL_WINDOW_MOUSE_FOCUS
	WindowFlagExternal          WindowFlags = C.SDL_WINDOW_EXTERNAL
	WindowFlagModal             WindowFlags = C.SDL_WINDOW_MODAL
	WindowFlagHighPixelDensity  WindowFlags = C.SDL_WINDOW_HIGH_PIXEL_DENSITY
	WindowFlagMouseCapture      WindowFlags = C.SDL_WINDOW_MOUSE_CAPTURE
	WindowFlagMouseRelativeMode WindowFlags = C.SDL_WINDOW_MOUSE_RELATIVE_MODE
	WindowFlagAlwaysOnTop       WindowFlags = C.SDL_WINDOW_ALWAYS_ON_TOP
	WindowFlagUtility           WindowFlags = C.SDL_WINDOW_UTILITY
	WindowFlagTooltip           WindowFlags = C.SDL_WINDOW_TOOLTIP
	WindowFlagPopupMenu         WindowFlags = C.SDL_WINDOW_POPUP_MENU
	WindowFlagKeyboardGrabbed   WindowFlags = C.SDL_WINDOW_KEYBOARD_GRABBED
	WindowFlagFillDocument      WindowFlags = C.SDL_WINDOW_FILL_DOCUMENT
	WindowFlagVulkan            WindowFlags = C.SDL_WINDOW_VULKAN
	WindowFlagMetal             WindowFlags = C.SDL_WINDOW_METAL
	WindowFlagTransparent       WindowFlags = C.SDL_WINDOW_TRANSPARENT
	WindowFlagNotFocusable      WindowFlags = C.SDL_WINDOW_NOT_FOCUSABLE
)

// Window is the opaque handle to an SDL window.
type Window struct {
	cWindow       *C.SDL_Window
	hitTestHandle cgo.Handle
}

// CreateWindow creates a window with the specified dimensions and flags.
//
// The window size is a request and may be different than expected based on
// desktop layout and window manager policies.
func CreateWindow(title string, width, height int, flags WindowFlags) (*Window, error) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	cw := C.SDL_CreateWindow(cTitle, C.int(width), C.int(height), C.Uint64(flags))
	if cw == nil {
		return nil, GetError()
	}

	return &Window{
		cWindow: cw,
	}, nil
}

// Destroy destroys a window.
func (w *Window) Destroy() {
	if w.hitTestHandle != cgo.Handle(0) {
		w.hitTestHandle.Delete()
		w.hitTestHandle = cgo.Handle(0)
	}
	C.SDL_DestroyWindow(w.cWindow)
	w.cWindow = nil
}

// CreateWindowAndRenderer creates a window and default renderer.
func CreateWindowAndRenderer(title string, width, height int, windowFlags WindowFlags) (*Window, Renderer, error) {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	cw := Window{}
	cr := cRenderer{}
	ok := C.SDL_CreateWindowAndRenderer(cTitle, C.int(width), C.int(height), C.Uint64(windowFlags), &cw.cWindow, &cr.cRenderer)
	if !ok {
		return nil, nil, GetError()
	}

	return &cw, &cr, nil
}

// SetWindowFullscreenMode sets the display mode to use when a window is visible and fullscreen.
func (w *Window) SetWindowFullscreenMode(mode *DisplayMode) error {
	var cMode *C.SDL_DisplayMode
	if mode != nil {
		cm := C.SDL_DisplayMode{
			displayID:                C.Uint32(mode.DisplayID),
			format:                   C.SDL_PixelFormat(mode.Format),
			w:                        C.int(mode.W),
			h:                        C.int(mode.H),
			pixel_density:            C.float(mode.PixelDensity),
			refresh_rate:             C.float(mode.RefreshRate),
			refresh_rate_numerator:   C.int(mode.RefreshRateNumerator),
			refresh_rate_denominator: C.int(mode.RefreshRateDenominator),
		}
		cMode = &cm
	}
	if !C.SDL_SetWindowFullscreenMode(w.cWindow, cMode) {
		return GetError()
	}
	return nil
}

// SetWindowFullscreen requests that the window's fullscreen state be changed.
func (w *Window) SetWindowFullscreen(fullscreen bool) error {
	if !C.SDL_SetWindowFullscreen(w.cWindow, C.bool(fullscreen)) {
		return GetError()
	}
	return nil
}

// GetWindowFlags gets the window flags.
func (w *Window) GetWindowFlags() WindowFlags {
	return WindowFlags(C.SDL_GetWindowFlags(w.cWindow))
}

// SetWindowTitle sets the title of a window.
//
// This string is expected to be in UTF-8 encoding.
func (w *Window) SetWindowTitle(title string) error {
	cTitle := C.CString(title)
	defer C.free(unsafe.Pointer(cTitle))
	if !C.SDL_SetWindowTitle(w.cWindow, cTitle) {
		return GetError()
	}
	return nil
}

// GetWindowTitle gets the title of a window.
func (w *Window) GetWindowTitle() string {
	title := C.SDL_GetWindowTitle(w.cWindow)
	if title == nil {
		return ""
	}
	return C.GoString(title)
}

// SetWindowIcon sets the icon for a window.
func (w *Window) SetWindowIcon(surface unsafe.Pointer) {
	C.SDL_SetWindowIcon(w.cWindow, (*C.SDL_Surface)(surface))
}

// SetWindowPosition requests that the window's position be set.
func (w *Window) SetWindowPosition(x, y int) {
	C.SDL_SetWindowPosition(w.cWindow, C.int(x), C.int(y))
}

// GetWindowPosition gets the position of a window.
func (w *Window) GetWindowPosition() (x, y int) {
	var cX, cY C.int
	C.SDL_GetWindowPosition(w.cWindow, &cX, &cY)
	return int(cX), int(cY)
}

// SetWindowSize requests that the size of a window's client area be set.
func (w *Window) SetWindowSize(width, height int) error {
	if !C.SDL_SetWindowSize(w.cWindow, C.int(width), C.int(height)) {
		return GetError()
	}
	return nil
}

// GetWindowSize gets the size of a window's client area.
func (w *Window) GetWindowSize() (width, height int) {
	var cW, cH C.int
	C.SDL_GetWindowSize(w.cWindow, &cW, &cH)
	return int(cW), int(cH)
}

// GetWindowSafeArea gets the safe area for this window.
func (w *Window) GetWindowSafeArea() (Rect, error) {
	var cRect C.SDL_Rect
	if !C.SDL_GetWindowSafeArea(w.cWindow, &cRect) {
		return Rect{}, GetError()
	}
	return Rect{X: int(cRect.x), Y: int(cRect.y), W: int(cRect.w), H: int(cRect.h)}, nil
}

// SetWindowBordered sets the border state of a window.
func (w *Window) SetWindowBordered(bordered bool) error {
	if !C.SDL_SetWindowBordered(w.cWindow, C.bool(bordered)) {
		return GetError()
	}
	return nil
}

// SetWindowResizable sets the user-resizable state of a window.
func (w *Window) SetWindowResizable(resizable bool) error {
	if !C.SDL_SetWindowResizable(w.cWindow, C.bool(resizable)) {
		return GetError()
	}
	return nil
}

// SetWindowAlwaysOnTop sets the window to always be above the others.
func (w *Window) SetWindowAlwaysOnTop(alwaysOnTop bool) error {
	if !C.SDL_SetWindowAlwaysOnTop(w.cWindow, C.bool(alwaysOnTop)) {
		return GetError()
	}
	return nil
}

// SetWindowFillDocument sets the window to fill the current document space (Emscripten only).
func (w *Window) SetWindowFillDocument(fillDocument bool) error {
	if !C.SDL_SetWindowFillDocument(w.cWindow, C.bool(fillDocument)) {
		return GetError()
	}
	return nil
}

// HideWindow hides a window.
func (w *Window) HideWindow() {
	C.SDL_HideWindow(w.cWindow)
}

// ShowWindow shows a window.
func (w *Window) ShowWindow() {
	C.SDL_ShowWindow(w.cWindow)
}

// RaiseWindow requests that a window be raised above other windows and gain input focus.
func (w *Window) RaiseWindow() {
	C.SDL_RaiseWindow(w.cWindow)
}

// MaximizeWindow requests that the window be made as large as possible.
func (w *Window) MaximizeWindow() {
	C.SDL_MaximizeWindow(w.cWindow)
}

// MinimizeWindow requests that the window be minimized to an iconic representation.
func (w *Window) MinimizeWindow() {
	C.SDL_MinimizeWindow(w.cWindow)
}

// RestoreWindow requests that the size and position of a minimized or maximized window be restored.
func (w *Window) RestoreWindow() {
	C.SDL_RestoreWindow(w.cWindow)
}

// SetWindowKeyboardGrab sets a window's keyboard grab mode.
func (w *Window) SetWindowKeyboardGrab(grabbed bool) error {
	if !C.SDL_SetWindowKeyboardGrab(w.cWindow, C.bool(grabbed)) {
		return GetError()
	}
	return nil
}

// SetWindowMouseGrab sets a window's mouse grab mode.
func (w *Window) SetWindowMouseGrab(grabbed bool) error {
	if !C.SDL_SetWindowMouseGrab(w.cWindow, C.bool(grabbed)) {
		return GetError()
	}
	return nil
}

// SetWindowRelativeMouseMode sets a window's relative mouse mode.
func (w *Window) SetWindowRelativeMouseMode(enabled bool) error {
	enabledInt := 0
	if enabled {
		enabledInt = 1
	}
	if C.sisrSetWindowRelativeMouseModeBridge(unsafe.Pointer(w.cWindow), C.int(enabledInt)) == 0 {
		return GetError()
	}
	return nil
}

// GetWindowKeyboardGrab gets a window's keyboard grab mode.
func (w *Window) GetWindowKeyboardGrab() bool {
	return bool(C.SDL_GetWindowKeyboardGrab(w.cWindow))
}

// GetWindowMouseGrab gets a window's mouse grab mode.
func (w *Window) GetWindowMouseGrab() bool {
	return bool(C.SDL_GetWindowMouseGrab(w.cWindow))
}

// GetWindowRelativeMouseMode gets a window's relative mouse mode.
func (w *Window) GetWindowRelativeMouseMode() bool {
	return C.sisrGetWindowRelativeMouseModeBridge(unsafe.Pointer(w.cWindow)) != 0
}

// SetWindowFocusable sets whether the window may have input focus.
func (w *Window) SetWindowFocusable(focusable bool) error {
	if !C.SDL_SetWindowFocusable(w.cWindow, C.bool(focusable)) {
		return GetError()
	}
	return nil
}

// SetWindowHitTest provides a callback that decides if a window region has special properties.
func (w *Window) SetWindowHitTest(callback WindowHitTestFunc) error {
	if callback == nil {
		if !C.SDL_SetWindowHitTest(w.cWindow, nil, nil) {
			return GetError()
		}
		if w.hitTestHandle != cgo.Handle(0) {
			w.hitTestHandle.Delete()
			w.hitTestHandle = cgo.Handle(0)
		}
		return nil
	}

	h := cgo.NewHandle(callback)
	if !C.SDL_SetWindowHitTest(
		w.cWindow,
		(C.SDL_HitTest)(C.sisrWindowHitTestBridge),
		unsafe.Pointer(uintptr(h)),
	) {
		h.Delete()
		return GetError()
	}

	if w.hitTestHandle != cgo.Handle(0) {
		w.hitTestHandle.Delete()
	}
	w.hitTestHandle = h

	return nil
}
