package sisr

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	apihandler "github.com/Alia5/SISR/api/handler"
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
	config.AutoUpdate             `embed:""`
	config.RunMode                `embed:""`
	config.ControllerEmulation    `embed:""`
	config.KeyboardMouseEmulation `embed:""`
	config.API                    `embed:"" prefix:"api."`
	config.Viiper                 `embed:"" prefix:"viiper."`
	config.Window                 `embed:"" prefix:"window."`
	config.Steam                  `embed:"" prefix:"steam."`
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
		err := steam.SetMarkerEnv(s.Steam.InstallDir, s.Steam.UserID) //nolint:staticcheck
		if err != nil {
			slog.Error("Failed to set Steam marker environment", "error", err)
		}
		slog.Info("Loading overlay...")
		err = steam.LoadOverlay(s.Steam.InstallDir)
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

	registerEventHandlers(eventRouter, &handler.Env{
		Window:          window,
		WebView:         wv,
		DeviceStore:     deviceStore,
		ViiperBridge:    viiperBridge,
		BindingEnforcer: bindingEnforcer,
		QuitFn:          stop,
		Config: &handler.RunConfig{
			AutoUpdate:             &s.AutoUpdate,
			RunMode:                &s.RunMode,
			ControllerEmulation:    &s.ControllerEmulation,
			KeyboardMouseEmulation: &s.KeyboardMouseEmulation,
			Viiper:                 &s.Viiper,
			Window:                 &s.Window,
			Steam:                  &s.Steam,
		},
	})

	_, apiAddr := s.runAPIServer(&apihandler.Env{
		Window:          window,
		WebView:         wv,
		DeviceStore:     deviceStore,
		BindingEnforcer: bindingEnforcer,
		QuitFn:          stop,
		Config: &apihandler.RunConfig{
			AutoUpdate:             &s.AutoUpdate,
			RunMode:                &s.RunMode,
			ControllerEmulation:    &s.ControllerEmulation,
			KeyboardMouseEmulation: &s.KeyboardMouseEmulation,
			Viiper:                 &s.Viiper,
			Window:                 &s.Window,
			Steam:                  &s.Steam,
		},
	})
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
