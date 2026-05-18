//go:build windows

package webview

import (
	"errors"

	"github.com/jchv/go-webview2/pkg/edge"

	"github.com/Alia5/SISR/sdl"
)

func Init() error {
	return nil
}

type windowsWebView struct {
	chromium *edge.Chromium
}

func New(window *sdl.Window, w, h int, debug bool) (WebView, error) {
	hwnd := window.GetPointerProperty(sdl.WindowPointerPropertyWin32HWND)
	if hwnd == 0 {
		return nil, errors.New("webview: could not get HWND from SDL window")
	}
	chromium := edge.NewChromium()
	if !chromium.Embed(hwnd) {
		return nil, errors.New("webview: failed to embed WebView2 into SDL window")
	}
	controller2 := chromium.GetController().GetICoreWebView2Controller2()
	if controller2 != nil {
		_ = controller2.PutDefaultBackgroundColor(edge.COREWEBVIEW2_COLOR{A: 0, R: 0, G: 0, B: 0})
	}
	chromium.Resize()
	return &windowsWebView{chromium: chromium}, nil
}

func (w *windowsWebView) Navigate(url string) {
	w.chromium.Navigate(url)
}

func (w *windowsWebView) SetHTML(html string) {
	w.chromium.NavigateToString(html)
}

func (w *windowsWebView) Eval(js string) {
	w.chromium.Eval(js)
}

func (w *windowsWebView) Bind(name string, fn interface{}) error {
	return errors.New("webview: Bind not implemented on Windows")
}

func (w *windowsWebView) Resize(width, height int) {
	w.chromium.Resize()
}

func (w *windowsWebView) Tick() {}

func (w *windowsWebView) Destroy() {}
