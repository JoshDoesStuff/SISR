//go:build linux

package webview

/*
#cgo pkg-config: webkitgtk-6.0 gtk4 x11
#include "webview_linux.h"
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"os"
	"unsafe"

	"github.com/Alia5/SISR/sdl"
)

func Init() error {
	os.Setenv("GDK_BACKEND", "x11")
	os.Setenv("SDL_VIDEO_DRIVER", "x11")
	os.Setenv("GSK_RENDERER", "cairo")
	os.Setenv("WEBKIT_DISABLE_COMPOSITING_MODE", "1")
	C.gtk_init_once()
	return nil
}

type linuxWebView struct {
	handle  *C.WebViewLinux
	visible bool
	width   int
	height  int
}

func New(window *sdl.Window, w, h int, _ bool) (WebView, error) {
	x11Window := window.GetX11WindowNumber()
	if x11Window == 0 {
		return nil, errors.New("webview: could not get X11 window ID from SDL window" +
			" – ensure webview.Init() was called before sdl.Init()")
	}
	handle := C.webview_create(C.ulong(x11Window), C.int(w), C.int(h))
	if handle == nil {
		return nil, errors.New("webview: webview_create returned NULL")
	}
	return &linuxWebView{handle: handle, visible: true, width: w, height: h}, nil
}

func (w *linuxWebView) Navigate(url string) {
	cs := C.CString(url)
	defer C.free(unsafe.Pointer(cs))
	C.webview_navigate(w.handle, cs)
}

func (w *linuxWebView) SetHTML(html string) {
	cs := C.CString(html)
	defer C.free(unsafe.Pointer(cs))
	C.webview_set_html(w.handle, cs)
}

func (w *linuxWebView) Eval(js string) {
	cs := C.CString(js)
	defer C.free(unsafe.Pointer(cs))
	C.webview_eval(w.handle, cs)
}

func (w *linuxWebView) Bind(name string, fn interface{}) error {
	return errors.New("webview: Bind not yet implemented on Linux")
}

func (w *linuxWebView) SetVisible(visible bool) {
	if visible {
		C.webview_set_visible(w.handle, C.int(1))
		C.webview_resize(w.handle, C.int(w.width), C.int(w.height))
	} else {
		C.webview_set_visible(w.handle, C.int(0))
	}
	w.visible = visible
}

func (w *linuxWebView) Visible() bool {
	return w.visible
}

func (w *linuxWebView) Resize(width, height int) {
	w.width = width
	w.height = height
	if w.visible {
		C.webview_resize(w.handle, C.int(width), C.int(height))
	}
}

func (w *linuxWebView) Tick() {
	C.webview_tick()
}

func (w *linuxWebView) Destroy() {
	if w.handle != nil {
		C.webview_destroy(w.handle)
		w.handle = nil
	}
}
