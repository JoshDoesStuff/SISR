package handler

import (
	"context"

	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/webview"
)

type Env struct {
	Window      *sdl.Window
	WebView     webview.WebView
	DeviceStore input.DeviceStore
	QuitFn      context.CancelFunc
}
