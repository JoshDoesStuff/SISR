package initiallaunch

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Alia5/SISR/cmd"
	"github.com/danielgtaylor/huma/v2"
)

type InitialLaunchResponse struct {
	Body InitialLaunch
}

type InitialLaunch struct {
	IsInitialLaunch bool `json:"is_initial_launch"`
}

func Register(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method: http.MethodGet,
		Path:   "/api/v1/initial_launch",
	}, status(c))
	huma.Register(a, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/api/v1/initial_launch",
		Description: "Writes the initial launch done file marker",
	}, writeInitialLaunchDoneFile(c))
}

func status(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*InitialLaunchResponse, error) {
	return func(ctx context.Context, req *struct{}) (*InitialLaunchResponse, error) {
		c.Config.Lock()
		defer c.Config.Unlock()

		return &InitialLaunchResponse{
			Body: InitialLaunch{
				IsInitialLaunch: c.Config.InitialLaunch,
			},
		}, nil
	}
}

func writeInitialLaunchDoneFile(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*InitialLaunchResponse, error) {
	return func(ctx context.Context, req *struct{}) (*InitialLaunchResponse, error) {

		ownExeDir, err := os.Executable()
		if err != nil {
			slog.Error("Failed to get own executable path", "error", err)
			return nil, err
		}
		ownExeDir, err = filepath.EvalSymlinks(ownExeDir)
		if err != nil {
			slog.Error("Failed to evaluate symlinks for own executable path", "error", err)
			return nil, err
		}

		markerPath := filepath.Join(filepath.Dir(ownExeDir), ".initial_setup_done")
		err = os.WriteFile(markerPath, []byte{}, 0644)
		if err != nil {
			slog.Error("Failed to write initial setup marker file", "error", err)
			return nil, err
		}

		c.Config.Lock()
		defer c.Config.Unlock()
		c.Config.InitialLaunch = false

		return &InitialLaunchResponse{
			Body: InitialLaunch{
				IsInitialLaunch: false,
			},
		}, nil
	}
}
