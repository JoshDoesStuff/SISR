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
	"github.com/Alia5/SISR/event"
	"github.com/Alia5/SISR/event/handler"
	"github.com/Alia5/SISR/helper"
	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/sdl/extras"
	"github.com/Alia5/SISR/webview"
)

type SISR struct {
	config.API   `embed:"" prefix:""`
	UpdateNotify config.UpdateNotify `help:"Update notification level: none, stable, prerelease" default:"stable" env:"SISR_UPDATE_NOTIFY"`

	Steam  config.Steam `embed:"" prefix:"steam."`
	MaxFPS uint32       `default:"60" help:"Maximim FPS for SteamOverlay/UI (Does not affect inputs)" env:"SISR_MAX_FPS"`
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

	if s.MaxFPS == 0 {
		s.targetFrameDuration = 0
	} else {
		s.targetFrameDuration = time.Second / 60
	}

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

	router := event.NewRouter()
	dh, dhCleanup, err := input.NewDeviceHandler(window, wv)
	if err != nil {
		slog.Error("Failed to initialite DeviceHanlder", "error", err)
	}
	defer dhCleanup()

	handlerParams := &handler.RegisterParams{
		Window:        window,
		WebView:       wv,
		DeviceHandler: dh,
		QuitFn:        stop,
	}
	registerEventHandlers(router, handlerParams)

	_, apiAddr := s.runAPIServer(window, wv, dh, stop)
	frontendAddr := s.FrontendAddress
	if frontendAddr == "" {
		frontendAddr = apiAddr
	}

	wv.Navigate(frontendAddr)

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

func registerEventHandlers(r event.Router, rp *handler.RegisterParams) {
	if runtime.GOOS == "linux" {
		hittestfunc := handler.HandleFunc(
			func(_ context.Context, ev *sdl.WindowEvent) error {
				return extras.HandleCursorHitTestWindowEvent(rp.Window, ev)
			},
		)
		event.RegisterHandler(r, handler.Operation[*sdl.WindowEvent]{
			Event:   sdl.EventTypeWindowPixelSizeChanged,
			Handler: hittestfunc,
		})
		event.RegisterHandler(r, handler.Operation[*sdl.WindowEvent]{
			Event:   sdl.EventTypeWindowResized,
			Handler: hittestfunc,
		})
	}
	event.RegisterHandler(r, handler.Quit(rp))
	event.RegisterHandler(r, handler.WindowResize(rp))
	event.RegisterHandler(r, handler.GamepadAdded(rp))
	event.RegisterHandler(r, handler.GamepadRemoved(rp))
	event.RegisterHandler(r, handler.GamepadUpdated(rp))
}

func cleanup() {
	slog.Info("Shutting down")
	err := helper.OpenURL("steam://forceinputappid/0")
	if err != nil {
		slog.Error("Failed to reset steam controller config", "error", err)
	}
}
