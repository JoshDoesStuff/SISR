package handler

import (
	"context"

	"github.com/Alia5/SISR/config"
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
	Config          *RunConfig
}

type RunConfig struct {
	*config.AutoUpdate
	*config.RunMode
	*config.ControllerEmulation
	*config.KeyboardMouseEmulation
	*config.Viiper
	*config.Window
	*config.Steam
}
