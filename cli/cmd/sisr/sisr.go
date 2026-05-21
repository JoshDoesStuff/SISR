package sisr

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alia5/SISR/config"
	"github.com/Alia5/SISR/event"
	"github.com/Alia5/SISR/helper"
	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/webview"
)

type SISR struct {
	config.API `embed:"" prefix:""`
	Steam      config.Steam `embed:"" prefix:"steam."`
	MaxFPS     uint32       `default:"60" help:"Maximim FPS for SteamOverlay/UI (Does not affect inputs)" env:"SISR_MAX_FPS"`
	//
	lastRenderTime      time.Time
	targetFrameDuration time.Duration
}

func (s *SISR) Run(cfg config.Global) error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt, syscall.SIGTERM,
	)
	defer stop()
	defer cleanup()

	setSDLHintEnv()
	setSDLHints()

	window, renderer, wv, err := s.createWindow(&cfg)
	if err != nil {
		return err
	}
	defer func() {
		wv.Destroy()
		renderer.Destroy()
		window.Destroy()
	}()

	dh, dhCleanup, err := input.NewDeviceHandler(window, wv)
	if err != nil {
		slog.Error("Failed to initialite DeviceHanlder", "error", err)
	}
	defer dhCleanup()

	_, apiAddr := s.runAPIServer(window, wv, dh, stop)
	frontendAddr := s.FrontendAddress
	if frontendAddr == "" {
		frontendAddr = apiAddr
	}
	router := event.NewRouter(window, renderer, wv, stop)
	wv.Navigate(frontendAddr)

	if s.MaxFPS == 0 {
		s.targetFrameDuration = 0
	} else {
		s.targetFrameDuration = time.Second / 60
	}

	return s.run(ctx, renderer, wv, router, stop)
}

func (s *SISR) run(
	ctx context.Context,
	renderer sdl.Renderer,
	wv webview.WebView,
	router event.Router,
	stop context.CancelFunc,
) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		ev, _ := sdl.WaitEventTimeout(s.targetFrameDuration)
		if ev != nil {
			router.RouteEvent(ctx, ev)
			// switch ev := ev.(type) {
			// case *sdl.QuitEvent:
			// 	return nil
			// case *sdl.KeyboardEvent:
			// 	if ev.Key == sdl.KeyCodeEscape && ev.Down {
			// 		wv.SetVisible(!wv.Visible())
			// 		if wv.Visible() {
			// 			slog.Info("WebView shown")
			// 		} else {
			// 			slog.Info("WebView hidden")
			// 		}
			// 	}
			// case *sdl.GamepadDeviceEvent:
			// 	switch ev.Type {
			// 	case sdl.EventTypeGamepadAdded:
			// 		id := sdl.GamepadID(ev.Which)
			// 		slog.Info("Gamepad connected", "id", id, "name", sdl.GetGamepadNameForID(id))
			// 		gp, err := sdl.OpenGamepad(id)
			// 		if err != nil {
			// 			slog.Error("Failed to open gamepad", "id", id, "error", err)
			// 		} else {
			// 			gps[id] = gp
			// 		}
			// 	case sdl.EventTypeGamepadRemoved:
			// 		id := sdl.GamepadID(ev.Which)
			// 		slog.Info("Gamepad disconnected", "id", id)
			// 		if gp, ok := gps[id]; ok {
			// 			gp.Close()
			// 			delete(gps, id)
			// 		}
			// 	}
			// case *sdl.WindowEvent:
			// 	if ev.Type == sdl.EventTypeWindowResized {
			// 		wv.Resize(int(ev.Data1), int(ev.Data2))
			// 	}
			// }
		}

		if time.Since(s.lastRenderTime) >= s.targetFrameDuration || ev == nil {
			s.lastRenderTime = time.Now()
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
}

func cleanup() {
	slog.Info("Shutting down")
	err := helper.OpenURL("steam://forceinputappid/0")
	if err != nil {
		slog.Error("Failed to reset steam controller config", "error", err)
	}
}
