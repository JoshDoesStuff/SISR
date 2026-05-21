package sisr

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/Alia5/SISR/config"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/sdl/extras"
	"github.com/Alia5/SISR/webview"
)

type SISR struct {
	config.API `embed:"" prefix:""`
	Steam      config.Steam `embed:"" prefix:"steam."`
}

func (s *SISR) Run(cfg config.Global) error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGTERM,
	)
	defer stop()

	defer func() {
		slog.Info("Shutting down")

	}()

	_, _, apiAddr := s.runAPIServer()
	frontendAddr := s.FrontendAddress
	if frontendAddr == "" {
		frontendAddr = apiAddr
	}

	window, renderer, wv, err := s.createWindow(&cfg, frontendAddr)
	if err != nil {
		return err
	}
	defer func() {
		wv.Destroy()
		renderer.Destroy()
		window.Destroy()
	}()

	sdl.SetGamepadEventsEnabled(true)

	return s.run(ctx, window, renderer, wv)
}

func (s *SISR) run(ctx context.Context, window *sdl.Window, renderer sdl.Renderer, wv webview.WebView) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		ev, _ := sdl.WaitEventTimeout(time.Millisecond * 16)
		if ev != nil {
			if runtime.GOOS == "linux" {
				err := extras.HandleCursorHitTestWindowEvent(window, ev)
				if err != nil {
					slog.Error("Failed to handle cursor hit test window event", "error", err)
				}
			}
			switch ev := ev.(type) {
			case *sdl.QuitEvent:
				return nil
			case *sdl.KeyboardEvent:
				if ev.Key == sdl.KeyCodeEscape && ev.Down {
					wv.SetVisible(!wv.Visible())
					if wv.Visible() {
						slog.Info("WebView shown")
					} else {
						slog.Info("WebView hidden")
					}
				}
			// case *sdl.GamepadDeviceEvent:
			// 	switch ev.Type {
			// 	case sdl.EventTypeGamepadAdded:
			// 		id := sdl.GamepadID(ev.Which)
			// 		slog.Info("Gamepad connected", "id", id, "name", sdl.GetGamepadNameForID(id))
			// 		openGamepad(id)
			// 	case sdl.EventTypeGamepadRemoved:
			// 		id := sdl.GamepadID(ev.Which)
			// 		slog.Info("Gamepad disconnected", "id", id)
			// 		closeGamepad(id)
			// 	}
			case *sdl.WindowEvent:
				if ev.Type == sdl.EventTypeWindowResized {
					wv.Resize(int(ev.Data1), int(ev.Data2))
				}
			}
		}

		wv.Tick()
		err := renderer.RenderClear()
		if err != nil {
			slog.Error("Failed to clear renderer", "error", err)
			return err
		}
		err = renderer.RenderPresent()
		if err != nil {
			slog.Error("Failed to present renderer", "error", err)
			return err
		}
	}
}
