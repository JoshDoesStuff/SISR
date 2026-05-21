package sisr

import (
	"log/slog"
	"runtime"

	"github.com/Alia5/SISR/config"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/webview"
	"github.com/Alia5/SISR/windows"
)

func (s *SISR) createWindow(
	cfg config.Global,
	frontendAddr string,
) (*sdl.Window, sdl.Renderer, webview.WebView, error) {
	setSDLHintEnv()
	err := sdl.Init(sdl.InitFlagVideo | sdl.InitFlagGamepad)
	if err != nil {
		slog.Error("Failed to init SDL", "error", err)
		return nil, nil, nil, err
	}

	flags := sdl.WindowFlagVulkan | sdl.WindowFlagTransparent | sdl.WindowFlagBorderless | sdl.WindowFlagAlwaysOnTop

	window, renderer, err := sdl.CreateWindowAndRenderer(
		"SISR",
		1280,
		720,
		flags,
	)
	if err != nil {
		slog.Error("Failed to create window", "error", err)
		sdl.Quit()
		return nil, nil, nil, err
	}

	if runtime.GOOS == "windows" {
		err = windows.SetDWMPassiveUpdateMode(window)
		if err != nil {
			slog.Error("Failed to enable DWM passive update mode", "error", err)
		}
	}

	err = renderer.SetRenderDrawColor(0, 0, 0, 0)
	if err != nil {
		slog.Error("Failed to set render draw color", "error", err)
		sdl.Quit()
		window.Destroy()
		return nil, nil, nil, err
	}

	wv, err := webview.New(window, 1280, 720, true)
	if err != nil {
		slog.Error("Failed to create webview", "error", err)
		sdl.Quit()
		window.Destroy()
		return nil, nil, nil, err
	}
	wv.SetHTML(`<!DOCTYPE html>
	<html style="background: transparent;">
	<body style="margin: 0; background: transparent; display: grid; place-items: center; height: 100svh;">
		<h1 style="color: red;">
			hello SISR ✂️
		</h1>
	</body>
	</html>`)
	// TODO: replace setHTML with navigate
	// wv.Navigate(frontendAddr)

	window.HideWindow()

	return window, renderer, wv, nil
}
