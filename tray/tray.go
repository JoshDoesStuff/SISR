package tray

import (
	"context"
	"fmt"
	"runtime"

	"fyne.io/systray"
	"github.com/Alia5/SISR/assets"
	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/meta"
)

type tray struct {
	*cmd.SISRContext
	exitItem *systray.MenuItem
}

func Run(ctx context.Context, env *cmd.SISRContext) {
	start, end := systray.RunWithExternalLoop(
		func() {
			onReady(ctx, env)
		},
		onExit,
	)
	start()
	go func() {
		<-ctx.Done()
		end()
	}()
}

func (t *tray) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			systray.Quit()
			return
		case <-t.exitItem.ClickedCh:
			t.QuitFn()
		}
	}
}

func onReady(ctx context.Context, env *cmd.SISRContext) {
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

	exitItem := systray.AddMenuItem("Quit", "Exit SISR")

	t := &tray{
		SISRContext: env,
		exitItem:    exitItem,
	}

	t.run(ctx)
}

func onExit() {
}
