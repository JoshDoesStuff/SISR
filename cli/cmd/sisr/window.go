package sisr

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/Alia5/SISR/config"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/sdl/extras"
	"github.com/Alia5/SISR/webview"
	"github.com/Alia5/SISR/windows"
)

func (s *SISR) createWindow(
	cfg *config.Window,
) (*sdl.Window, sdl.Renderer, webview.WebView, error) {

	err := sdl.Init(sdl.InitFlagVideo)
	if err != nil {
		slog.Error("Failed to init SDL", "error", err)
		return nil, nil, nil, err
	}

	flags := sdl.WindowFlagVulkan | sdl.WindowFlagTransparent
	if cfg.Fullscreen {
		flags |= sdl.WindowFlagBorderless | sdl.WindowFlagAlwaysOnTop | sdl.WindowFlagFullscreen
	} else {
		flags |= sdl.WindowFlagResizable
	}

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

	wv, err := webview.New(window, 1280, 720, os.Getenv("DEV") == "1")
	if err != nil {
		slog.Error("Failed to create webview", "error", err)
		sdl.Quit()
		window.Destroy()
		return nil, nil, nil, err
	}

	if !cfg.Show {
		window.HideWindow()
	}
	if cfg.Fullscreen {
		wv.SetVisible(false)
		err := extras.SetCursorHitTest(window, false)
		if err != nil {
			slog.Error("Failed setting window cursor hittest", "error", err)
			err = window.SetWindowFullscreen(false)
			if err != nil {
				slog.Error("Failed to set window fullscreen", "error", err)
			}
			err = window.SetWindowResizable(true)
			if err != nil {
				slog.Error("Failed to set window resizable", "error", err)
			}
		}
	}

	return window, renderer, wv, nil
}
