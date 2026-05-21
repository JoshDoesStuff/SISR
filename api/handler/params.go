package handler

import (
	"context"

	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/webview"
)

type RegisterParams struct {
	Window        *sdl.Window
	WebView       webview.WebView
	DeviceHandler input.DeviceHandler
	QuitFn        context.CancelFunc
}
