package input

import (
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/webview"
)

type DeviceHandler interface {
}

type deviceHandler struct {
	window   *sdl.Window
	webview  webview.WebView
	gamepads map[sdl.GamepadID]*sdl.Gamepad
}

func NewDeviceHandler(window *sdl.Window, webview webview.WebView) (DeviceHandler, func(), error) {
	sdl.SetGamepadEventsEnabled(true)
	err := sdl.InitSubSystem(sdl.InitFlagGamepad | sdl.InitFlagSensor | sdl.InitFlagHaptic)
	if err != nil {
		return nil, nil, err
	}

	dh := &deviceHandler{
		window:   window,
		webview:  webview,
		gamepads: make(map[sdl.GamepadID]*sdl.Gamepad),
	}
	return dh, dh.quit, nil
}

func (dh *deviceHandler) quit() {
	for id, gp := range dh.gamepads {
		if gp != nil {
			gp.Close()
		}
		delete(dh.gamepads, id)
	}
}
