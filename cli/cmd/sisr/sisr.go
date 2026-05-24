package sisr

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/config"
	"github.com/Alia5/SISR/event"
	"github.com/Alia5/SISR/helper"
	"github.com/Alia5/SISR/input"
	"github.com/Alia5/SISR/input/steaminputbindings"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/steam"
	"github.com/Alia5/SISR/tray"
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
		err = steam.LoadOverlay(s.Steam.InstallDir) //nolint:staticcheck
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
	viiperBridge := input.NewViiperBridge(ctx, deviceStore, &s.Viiper)

	winDispatcher := cmd.NewWindowDispatcher[any]()
	cmdCtx := &cmd.SISRContext{
		WindowDispatcher: winDispatcher,
		DeviceStore:      deviceStore,
		ViiperBridge:     viiperBridge,
		BindingEnforcer:  bindingEnforcer,
		QuitFn:           stop,
		Config: &cmd.SessionConfig{
			AutoUpdate:             &s.AutoUpdate,
			RunMode:                &s.RunMode,
			ControllerEmulation:    &s.ControllerEmulation,
			KeyboardMouseEmulation: &s.KeyboardMouseEmulation,
			Viiper:                 &s.Viiper,
			Window:                 &s.Window,
			Steam:                  &s.Steam,
		},
	}

	registerEventHandlers(eventRouter, cmdCtx, window, wv)

	_, apiAddr := s.runAPIServer(cmdCtx)
	frontendAddr := s.FrontendAddress
	if frontendAddr == "" {
		frontendAddr = apiAddr
	}

	wv.Navigate(frontendAddr)

	// TODO: check settings and stuff
	err = bindingEnforcer.ForceOwnAppID()
	if err != nil {
		slog.Error("Failed to force SteamInput layout", "error", err)
	}

	return s.run(ctx, renderer, window, wv, eventRouter, winDispatcher)
}

func (s *SISR) run(
	ctx context.Context,
	renderer sdl.Renderer,
	window *sdl.Window,
	wv webview.WebView,
	router event.Router,
	dispatcher cmd.WindowDispatcher[any],
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
			drainEvents(ctx, router)
		}
		dispatcher.Dispatch(window, wv)
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

func drainEvents(ctx context.Context, router event.Router) {
	for {
		ev, _ := sdl.PollEvent()
		if ev == nil {
			break
		}
		router.RouteEvent(ctx, ev)
	}
}

func cleanup() {
	slog.Info("Shutting down")
	err := helper.OpenURL("steam://forceinputappid/0")
	if err != nil {
		slog.Error("Failed to reset steam controller config", "error", err)
	}
}
