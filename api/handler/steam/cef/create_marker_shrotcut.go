package cef

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Alia5/SISR/cefpayloads"
	"github.com/Alia5/SISR/cmd"
	"github.com/danielgtaylor/huma/v2"
)

func registerCreateMarkerShortcut(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/steam/cef/create-marker-shortcut",
		Tags:   []string{"steam", "cef", "inject"},
	}, createMarkerShortcut(c))
}

func createMarkerShortcut(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*struct{}, error) {
	return func(ctx context.Context, req *struct{}) (*struct{}, error) {
		c.Config.Lock()
		steamCfg := *c.Config.Steam
		c.Config.Unlock()

		SISRPath, err := os.Executable()
		if err != nil {
			slog.Error("Failed to get SISR executable path", "error", err)
			return nil, err
		}

		steamAppID := os.Getenv("SteamAppId")
		markerAppID := os.Getenv("SISR_MARKER_ID")
		if (steamAppID != "" && steamAppID != "0") || (markerAppID != "" && markerAppID != "0") {
			slog.Debug("Marker shortcut already exists OR launched via Steam, skipping creation", "appId", markerAppID)
			return nil, nil
		}

		_, err = cefpayloads.NewCreateMarkerShortcut(&steamCfg).
			Execute(
				ctx,
				&cefpayloads.CreateMarkerShortcutArgs{
					SISRPath: filepath.ToSlash(SISRPath),
				},
			)
		if err != nil {
			slog.Error("failed to create marker shortcut", "error", err)
		}

		return nil, err
	}
}
