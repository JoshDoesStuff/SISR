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

	allowDesktopLayoutItem := systray.AddMenuItemCheckbox(
		"Allow Steam Desktop Layout",
		"Uses Steams Desktop Layout instead of SISR Marker (or specific SISR layout when launched via Steam)",
		c.Config.AllowSteamDesktopLayout,
	)

	systray.AddSeparator()

	exitItem := systray.AddMenuItem("Quit", "Exit SISR")

	t := &tray{
		SISRContext:            c,
		exitItem:               exitItem,
		allowDesktopLayoutItem: allowDesktopLayoutItem,
	}

	t.run(ctx)
}

type tray struct {
	*cmd.SISRContext
	exitItem               *systray.MenuItem
	allowDesktopLayoutItem *systray.MenuItem
}

func (t *tray) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			systray.Quit()
			return
		case <-t.exitItem.ClickedCh:
			t.QuitFn()
		case <-t.allowDesktopLayoutItem.ClickedCh:
			t.handleAllowDesktopConfig()
		}
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
