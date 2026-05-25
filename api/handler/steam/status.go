package steam

import (
	"context"
	"log/slog"
	"os"
	"strconv"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/steam"
)

type SteamStatusResponse struct {
	Body SteamAndCefStatus
}

type SteamAndCefStatus struct {
	SteamRunning      bool   `json:"steam_running"`
	SteamPath         string `json:"steam_path"`
	CEFDebugEnabled   bool   `json:"cef_debug_enabled"`
	CEFDebugReachable bool   `json:"cef_debug_reachable"`

	NoSteamMode bool `json:"no_steam_mode"`

	LaunchedViaSteam bool   `json:"launched_via_steam"`
	SteamGameID      string `json:"steam_game_id"`
	SteamAppID       uint32 `json:"steam_app_id"`
}

func status(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*SteamStatusResponse, error) {
	return func(ctx context.Context, req *struct{}) (*SteamStatusResponse, error) {

		c.Config.Lock()
		defer c.Config.Unlock()

		if c.Config.NoSteam {
			return &SteamStatusResponse{
				Body: SteamAndCefStatus{
					NoSteamMode: true,
				},
			}, nil
		}

		steamRunning := steam.ClientRunning()
		steamPath := c.Config.Steam.InstallDir // nolint
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
