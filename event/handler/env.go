package handler

import (
	"context"

	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/input/steaminputbindings"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/webview"
)

type Env struct {
	Window          *sdl.Window
	WebView         webview.WebView
	DeviceStore     input.DeviceStore
	ViiperBridge    input.ViiperBridge
	BindingEnforcer steaminputbindings.Enforcer
	QuitFn          context.CancelFunc
}
