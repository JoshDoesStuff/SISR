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
	"github.com/Alia5/SISR/event/handler"
	"github.com/Alia5/SISR/helper"
	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/input/steaminputbindings"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/steam"
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

	launchedViaSteam, launchedInGameMode := steam.LaunchedViaSteam()
	_ = launchedInGameMode

	if !launchedViaSteam {
		slog.Info("Not launched via Steam, setting env...")
		err := steam.SetMarkerEnv()
		if err != nil {
			slog.Error("Failed to set Steam marker environment", "error", err)
		}
		slog.Info("Loading overlay...")
		err = steam.LoadOverlay()
		if err != nil {
			slog.Error("Failed to load Steam overlay", "error", err)
		}
	} else {
		slog.Info("Launched via Steam")
	}

	window, renderer, wv, err := s.createWindow(&cfg)
	if err != nil {
		return err
	}
	defer func() {
		wv.Destroy()
		renderer.Destroy()
		window.Destroy()
	}()

	unhookSteamHid()
	err = sdl.InitSubSystem(sdl.InitFlagGamepad | sdl.InitFlagSensor | sdl.InitFlagHaptic)
	if err != nil {
		return err
	}

	bindingEnforcer := steaminputbindings.NewEnforcer()
	eventRouter := event.NewRouter()
	deviceStore, deviceStoreClose, err := input.NewDeviceStore()
	if err != nil {
		slog.Error("Failed to initialite DeviceStore", "error", err)
	}
	defer deviceStoreClose()
	viiperBridge := input.NewViiperBridge(ctx, deviceStore)

	handlerEnv := &handler.Env{
		Window:          window,
		WebView:         wv,
		DeviceStore:     deviceStore,
		ViiperBridge:    viiperBridge,
		BindingEnforcer: bindingEnforcer,
		QuitFn:          stop,
	}
	registerEventHandlers(eventRouter, handlerEnv)

	_, apiAddr := s.runAPIServer(window, wv, deviceStore, bindingEnforcer, stop)
	frontendAddr := s.FrontendAddress
	if frontendAddr == "" {
		frontendAddr = apiAddr
	}

	wv.Navigate(frontendAddr)

	return s.run(ctx, renderer, wv, eventRouter)
}

func (s *SISR) run(
	ctx context.Context,
	renderer sdl.Renderer,
	wv webview.WebView,
	router event.Router,
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

func cleanup() {
	slog.Info("Shutting down")
	err := helper.OpenURL("steam://forceinputappid/0")
	if err != nil {
		slog.Error("Failed to reset steam controller config", "error", err)
	}
}
