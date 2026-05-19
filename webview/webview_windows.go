//go:build windows

package webview

import (
	"errors"
	"log/slog"

	"github.com/jchv/go-webview2/pkg/edge"

	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/windows"
)

func Init() error {
	return nil
}

type windowsWebView struct {
	chromium *edge.Chromium
	visible  bool
	width    int
	height   int
}

func prepareWindowForWebView(hwnd uintptr) error {
	return windows.ClearWindowStyleBits(hwnd, windows.WSClipChildren)
}

func New(window *sdl.Window, w, h int, debug bool) (WebView, error) {
	hwnd := window.GetPointerProperty(sdl.WindowPointerPropertyWin32HWND)
	if hwnd == 0 {
		return nil, errors.New("webview: could not get HWND from SDL window")
	}
	if err := prepareWindowForWebView(hwnd); err != nil {
		return nil, err
	}
	chromium := edge.NewChromium()
	if !chromium.Embed(hwnd) {
		return nil, errors.New("webview: failed to embed WebView2 into SDL window")
	}
	controller2 := chromium.GetController().GetICoreWebView2Controller2()
	if controller2 != nil {
		err := controller2.PutDefaultBackgroundColor(edge.COREWEBVIEW2_COLOR{
			A: 0, R: 0, G: 0, B: 0,
		})
		if err != nil {
			slog.Error("Failed to set WebView2 default background color", "error", err)
		}
	}
	chromium.Resize()
	return &windowsWebView{chromium: chromium, visible: true, width: w, height: h}, nil
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

func (w *windowsWebView) Bind(name string, fn any) error {
	return errors.New("webview: Bind not implemented on Windows")
}

func (w *windowsWebView) SetVisible(visible bool) {
	if visible {
		err := w.chromium.Show()
		if err != nil {
			slog.Error("Failed to show WebView2", "error", err)
		}
		w.chromium.Resize()
	} else {
		err := w.chromium.Hide()
		if err != nil {
			slog.Error("Failed to hide WebView2", "error", err)
		}
	}
	w.visible = visible
}

func (w *windowsWebView) Visible() bool {
	return w.visible
}

func (w *windowsWebView) Resize(width, height int) {
	w.width = width
	w.height = height
	if w.visible {
		w.chromium.Resize()
	}
}

func (w *windowsWebView) Tick() {}

func (w *windowsWebView) Destroy() {}
