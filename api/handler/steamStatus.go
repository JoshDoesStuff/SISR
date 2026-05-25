package handler

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/steam"
	"github.com/danielgtaylor/huma/v2"
)

type SteamStatusResponse struct {
	Body SteamAndCefStatus
}

type SteamAndCefStatus struct {
	SteamRunning      bool   `json:"steam_running"`
	SteamPath         string `json:"steam_path"`
	CEFDebugEnabled   bool   `json:"cef_debug_enabled"`
	CEFDebugReachable bool   `json:"cef_debug_reachable"`

	LaunchedViaSteam bool   `json:"launched_via_steam"`
	SteamGameID      string `json:"steam_game_id"`
	SteamAppID       uint32 `json:"steam_app_id"`
}

func RegisterSteamStatus(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method: http.MethodGet,
		Path:   "/api/v1/steam_status",
	}, steamStatusHandleFunc(c))
}

func steamStatusHandleFunc(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*SteamStatusResponse, error) {
	return func(ctx context.Context, req *struct{}) (*SteamStatusResponse, error) {

		c.Config.Lock()
		defer c.Config.Unlock()

		steamRunning := steam.ClientRunning()
		steamPath := c.Config.Steam.InstallDir
		if steamPath == "" {
			var err error
			steamPath, err = steam.ExecuteableDir()
			if err != nil {
				slog.Error("Failed to detect Steam path", "error", err)
			}
		}
		debugEnableFilePresent, err := steam.CEFRemoteDebugEnableFilePresent(c.Config.Steam)
		if err != nil {
			slog.Error("Failed to check CEF debug enable file", "error", err)
		}
		debugReachable := false
		if steamRunning || debugEnableFilePresent {
			debugReachable = steam.CEFRemoteDebugReachable(ctx, c.Config.Steam)
		}
		if !steamRunning && debugReachable {
			steamRunning = true
		}
		if debugReachable {
			debugEnableFilePresent = true
		}
		launchedViaSteam, _ := steam.LaunchedViaSteam()
		steamAppIDStr := os.Getenv("SteamAppId")
		if steamAppIDStr == "" || steamAppIDStr == "0" {
			steamAppIDStr = os.Getenv("SISR_MARKER_ID")
		}
		steamAppID, err := strconv.ParseUint(steamAppIDStr, 10, 32)
		if err != nil {
			steamAppID = 0
			slog.Error("Failed to parse Steam App ID", "error", err, "value", steamAppIDStr)
		}

		steamGameID := os.Getenv("SteamGameId")

		return &SteamStatusResponse{
			Body: SteamAndCefStatus{
				SteamRunning:      steamRunning,
				SteamPath:         steamPath,
				CEFDebugEnabled:   debugEnableFilePresent,
				CEFDebugReachable: debugReachable,
				LaunchedViaSteam:  launchedViaSteam,
				SteamAppID:        uint32(steamAppID),
				SteamGameID:       steamGameID,
			},
		}, nil
	}
}
