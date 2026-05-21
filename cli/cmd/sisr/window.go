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
	cfg *config.Global, //nolint:unparam
) (*sdl.Window, sdl.Renderer, webview.WebView, error) {

	err := sdl.Init(sdl.InitFlagVideo)
	if err != nil {
		slog.Error("Failed to init SDL", "error", err)
		return nil, nil, nil, err
	}

	flags := sdl.WindowFlagVulkan | sdl.WindowFlagTransparent
	// sdl.WindowFlagBorderless | sdl.WindowFlagAlwaysOnTop // TODO:

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

	// window.HideWindow() // TODO:

	return window, renderer, wv, nil
}
