package steam

import (
	"errors"
	"log/slog"
	"os"
	"path"
	"runtime"

	"github.com/Alia5/SISR/helper"
)

var ErrOverlayLoadLaunchedViaSteam = errors.New("launched via Steam, overlay should already be loaded")
var ErrSteamNotRunning = errors.New("Steam is not running, loading overlay is useless")

func LoadOverlay() error {
	if launchedViaSteam, _ := LaunchedViaSteam(); launchedViaSteam {
		return ErrOverlayLoadLaunchedViaSteam
	}

	if !steamRunning() {
		return ErrSteamNotRunning
	}

	sp, err := steamPath()
	if err != nil {
		return err
	}

	var overlayPath string
	switch runtime.GOOS {
	case "windows":
		overlayPath = path.Join(sp, "GameOverlayRenderer64.dll")
	case "linux":
		parentDir := path.Dir(sp)
		ubuntu12_64 := path.Join(parentDir, "ubuntu12_64", "gameoverlayrenderer.so")
		bin64 := path.Join(parentDir, "bin64", "gameoverlayrenderer.so")

		if _, err := os.Stat(bin64); err == nil {
			overlayPath = bin64
		} else {
			overlayPath = ubuntu12_64
		}
	}
	slog.Debug("Attempting to load steamoverlay", "path", overlayPath)
	err = helper.LoadLib(overlayPath)
	if err != nil {
		return err
	}

	slog.Info("Successfully loaded steam overlay")

	return nil
}
