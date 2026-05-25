package tray

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"

	"fyne.io/systray"
	"github.com/Alia5/SISR/assets"
	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/meta"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/sdl/extras"
	"github.com/Alia5/SISR/webview"
)

func Run(ctx context.Context, c *cmd.SISRContext) {
	start, end := systray.RunWithExternalLoop(
		func() {
			onReady(ctx, c)
		},
		func() {},
	)
	start()
	go func() {
		<-ctx.Done()
		end()
	}()
}

func onReady(ctx context.Context, c *cmd.SISRContext) {
	if runtime.GOOS == "windows" {
		systray.SetIcon(assets.IconICO)
	} else {
		systray.SetIcon(assets.IconPNG)
	}
	systray.SetTooltip("SISR - Steam Input System Redirector")

	infoStr := fmt.Sprintf("SISR - %s", meta.Version)
	versionItem := systray.AddMenuItem(infoStr, infoStr)
	versionItem.Disable()

	systray.AddSeparator()

	toggleUIItem := systray.AddMenuItem("Toggle UI", "Show or hide the SISR UI")
	enableOverlayItem := systray.AddMenuItemCheckbox(
		"Enable Steam Overlay",
		"Enables the Steam overlay as transparent borderless window",
		c.Config.Fullscreen,
	)

	systray.AddSeparator()

	allowDesktopLayoutItem := systray.AddMenuItemCheckbox(
		"Allow Steam Desktop Layout",
		"Uses Steams Desktop Layout instead of SISR Marker (or specific SISR layout when launched via Steam)",
		c.Config.AllowSteamDesktopLayout,
	)

	systray.AddSeparator()

	exitItem := systray.AddMenuItem("Quit", "Exit SISR")

	t := &tray{
		SISRContext: c,

		toggleUIItem:      toggleUIItem,
		enableOverlayItem: enableOverlayItem,

		allowDesktopLayoutItem: allowDesktopLayoutItem,

		exitItem: exitItem,
	}

	t.run(ctx)
}

type tray struct {
	*cmd.SISRContext

	toggleUIItem      *systray.MenuItem
	enableOverlayItem *systray.MenuItem

	allowDesktopLayoutItem *systray.MenuItem

	exitItem *systray.MenuItem
}

func (t *tray) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			systray.Quit()
			return
		case <-t.toggleUIItem.ClickedCh:
			t.handleToggleUI(ctx)
		case <-t.enableOverlayItem.ClickedCh:
			t.handleToggleOverlay(ctx)
		case <-t.allowDesktopLayoutItem.ClickedCh:
			t.handleAllowDesktopConfig()
		case <-t.exitItem.ClickedCh:
			t.QuitFn()
		}
	}
}

func (t *tray) handleToggleUI(ctx context.Context) {
	_, err := cmd.ScheduleWindowDispatch(ctx, t.WindowDispatcher, func(w *sdl.Window, wv webview.WebView) bool {
		t.Config.Lock()
		fullscreen := t.Config.Fullscreen
		t.Config.Unlock()
		if wv.Visible() {
			if !fullscreen {
				w.HideWindow()
			}
			wv.SetVisible(false)
			err := extras.SetCursorHitTest(w, false)
			if err != nil {
				slog.Error("Failed setting window cursor hittest", "error", err)
			}
			return false
		} else {
			w.ShowWindow()
			// wv.SetVisible(true)
			_ = t.WindowDispatcher.Schedule(func(w *sdl.Window, wv webview.WebView) any {
				wv.SetVisible(true)
				return nil
			})
			err := extras.SetCursorHitTest(w, true)
			if err != nil {
				slog.Error("Failed setting window cursor hittest", "error", err)
			}
			return true
		}
	})
	if err != nil {
		slog.Error("Failed to toggle UI visibility", "error", err)
	}
}

func (t *tray) handleToggleOverlay(ctx context.Context) {
	if t.enableOverlayItem.Checked() {
		t.enableOverlayItem.Uncheck()
	} else {
		t.enableOverlayItem.Check()
	}
	t.Config.Lock()
	t.Config.Fullscreen = t.enableOverlayItem.Checked()
	fullscreen := t.Config.Fullscreen
	t.Config.Unlock()

	err, dispatchErr := cmd.ScheduleWindowDispatch(ctx, t.WindowDispatcher, func(w *sdl.Window, wv webview.WebView) error {
		if wv.Visible() {
			wv.SetVisible(false)
		}
		if fullscreen {
			w.ShowWindow()
			err := extras.SetCursorHitTest(w, false)
			if err != nil {
				slog.Error("Failed setting window cursor hittest", "error", err)
			}
			err = w.SetWindowFullscreen(true)
			if err != nil {
				slog.Debug("Failed to set window fullscreen", "error", err)
				return err
			}
			err = w.SetWindowAlwaysOnTop(true)
			if err != nil {
				slog.Debug("Failed to set window always on top", "error", err)
				return err
			}
			err = w.SetWindowResizable(false)
			if err != nil {
				slog.Debug("Failed to set window resizable", "error", err)
				return err
			}
			err = w.SetWindowBordered(false)
			if err != nil {
				slog.Debug("Failed to set window bordered", "error", err)
				return err
			}
		} else {
			w.HideWindow()
			_ = t.WindowDispatcher.Schedule(func(w *sdl.Window, wv webview.WebView) any {
				wv.SetVisible(true)
				return nil
			})
			err := extras.SetCursorHitTest(w, true)
			if err != nil {
				slog.Error("Failed setting window cursor hittest", "error", err)
			}
			err = w.SetWindowFullscreen(false)
			if err != nil {
				slog.Debug("Failed to set window fullscreen", "error", err)
				return err
			}
			err = w.SetWindowAlwaysOnTop(false)
			if err != nil {
				slog.Debug("Failed to set window always on top", "error", err)
				return err
			}
			err = w.SetWindowResizable(false)
			if err != nil {
				slog.Debug("Failed to set window resizable", "error", err)
				return err
			}
			err = w.SetWindowBordered(false)
			if err != nil {
				slog.Debug("Failed to set window bordered", "error", err)
				return err
			}
		}
		return nil
	})
	if dispatchErr != nil {
		slog.Error("Error dispatching fullscreen state change to window", "error", err)
	}
	if err != nil {
		slog.Error("Failed to change fullscreen state", "error", err)
	}

}

func (t *tray) handleAllowDesktopConfig() {
	if t.allowDesktopLayoutItem.Checked() {
		t.allowDesktopLayoutItem.Uncheck()
	} else {
		t.allowDesktopLayoutItem.Check()
	}
	if t.allowDesktopLayoutItem.Checked() {
		err := t.BindingEnforcer.ForceInputAppID(0) // 0 reverts to non-forced, including desktop-layout when out of focus
		if err != nil {
			slog.Error("Failed to reset forced inputAppID", "error", err)
			t.allowDesktopLayoutItem.Uncheck()
			return
		}
	} else {
		err := t.BindingEnforcer.ForceOwnAppID()
		if err != nil {
			slog.Error("Failed to force SteamInput layout", "error", err)
			t.allowDesktopLayoutItem.Check()
			return
		}
	}
	t.Config.Lock()
	t.Config.AllowSteamDesktopLayout = t.allowDesktopLayoutItem.Checked()
	t.Config.Unlock()
}
