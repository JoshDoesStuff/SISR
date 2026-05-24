package steam

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Alia5/SISR/steam/vdf"
)

func SetMarkerEnv(steamDir string, steamUserID uint32) error {

	if steamDir == "" {
		var err error
		steamDir, err = steamPath()
		if err != nil {
			slog.Error("Could not detect Steam Path", "error", err)
			return err
		}
	}

	if steamUserID == 0 {
		var err error
		steamUserID, err = ActiveUserID()
		if err != nil {
			slog.Error("Could not detect active Steam user ID", "error", err)
			return err
		}
	}

	shortcutsPath, err := ShortcutsPath(steamDir, steamUserID)
	if err != nil {
		slog.Error("Could not determine Steam shortcuts.vdf path", "error", err)
		return err
	}

	markerAppID, err := MarkerAppID(shortcutsPath)
	if err != nil {
		slog.Error("Failed to check for SISR marker in Steam shortcuts", "error", err)
		return err
	}
	if markerAppID == 0 {
		slog.Error("No SISR marker shortcut found in Steam shortcuts; Steam integration may not work correctly")
		return ErrMarkerNotFound
	}

	// avoid showing ingame status in steam
	_ = os.Setenv("SteamClientLaunch", "0")
	_ = os.Setenv("SteamAppId", "0")
	//
	_ = os.Setenv("SISR_MARKER_ID", strconv.FormatUint(uint64(markerAppID), 10))
	gameID := (uint64(markerAppID) << 32) | (2 << 24)
	_ = os.Setenv("SteamGameId", strconv.FormatUint(gameID, 10))
	_ = os.Setenv("SteamOverlayGameId", strconv.FormatUint(gameID, 10))
	_ = os.Setenv("SteamPath", steamDir)

	_ = os.Setenv("SteamVirtualGamepadInfo", path.Join(steamDir, "config", "virtualgamepadinfo.txt"))

	return nil
}

func MarkerAppID(shortcutsPath string) (uint32, error) {

	file, err := os.Open(shortcutsPath)
	if err != nil {
		return 0, err
	}
	defer file.Close() //nolint:errcheck

	vdfData, err := vdf.Read(file)
	if err != nil {
		return 0, err
	}

	runningPath := os.Getenv("APPIMAGE")
	if runningPath == "" {
		runningPath, err = os.Executable()
		if err != nil {
			return 0, err
		}
	}

	resolvedPath, err := filepath.EvalSymlinks(runningPath)
	if err == nil {
		runningPath = resolvedPath
	}
	absolutePath, err := filepath.Abs(runningPath)
	if err == nil {
		runningPath = absolutePath
	}
	runningPath = strings.ToLower(strings.ReplaceAll(runningPath, "\\", "/"))

	shortcutsValue, ok := vdfData["shortcuts"]
	if !ok {
		return 0, nil
	}

	shortcutsArray, ok := shortcutsValue.(map[string]any)
	if !ok {
		return 0, nil
	}

	for _, shortcut := range shortcutsArray {
		shortcutMap, ok := shortcut.(map[string]any)
		if !ok {
			continue
		}

		var pathValue any
		var argsValue any
		for key, value := range shortcutMap {
			if strings.EqualFold(key, "exe") {
				pathValue = value
			}

			if strings.EqualFold(key, "LaunchOptions") {
				argsValue = value
			}
		}

		pathStr, ok := pathValue.(string)
		if !ok {
			continue
		}

		argsStr, ok := argsValue.(string)
		if !ok {
			continue
		}

		if strings.Contains(
			strings.ToLower(strings.ReplaceAll(pathStr, "\\", "/")),
			runningPath,
		) && strings.Contains(
			strings.ToLower(argsStr),
			"--marker",
		) {
			var appIDValue any
			for key, value := range shortcutMap {
				if strings.EqualFold(key, "appid") {
					appIDValue = value
					break
				}
			}

			if appIDValue != nil {
				parsedID, parseErr := strconv.ParseUint(
					fmt.Sprint(appIDValue), 10, 32,
				)
				if parseErr == nil {
					return uint32(parsedID), nil
				}
			}

			return 0, nil
		}
	}

	return 0, nil

}

// let Some(steam_path) = steam_path() else {
//     warn!("Steam path could not be determined; Steam integration may not work correctly");
//     return Err(anyhow::anyhow!("Steam path could not be determined"));
// };
// let Some(steam_active_user_id) = active_user_id() else {
//     warn!(
//         "Active Steam user ID could not be determined; Steam integration may not work correctly"
//     );
//     return Err(anyhow::anyhow!(
//         "Active Steam user ID could not be determined"
//     ));
// };
// let Some(shortcuts_path) = get_shortcuts_path(&steam_path.clone(), steam_active_user_id) else {
//     warn!("Failed to determine Steam shortcuts.vdf path");
//     return Err(anyhow::anyhow!(
//         "Failed to determine Steam shortcuts.vdf path"
//     ));
// };
// trace!("Steam shortcuts.vdf path: {:?}", shortcuts_path);
// let marker_app_id = shortcuts_has_sisr_marker(&shortcuts_path);
// if marker_app_id == 0 {
//     warn!(
//         "No SISR marker shortcut found in Steam shortcuts; Steam integration may not work correctly"
//     );
//     return Err(anyhow::anyhow!(
//         "No SISR marker shortcut found in Steam shortcuts"
//     ));
// }
// unsafe {
//     std::env::set_var("SteamClientLaunch", "0");

//     std::env::set_var("SteamAppId", "0");
//     std::env::set_var("SISR_MARKER_ID", marker_app_id.to_string());
//     let game_id = (marker_app_id as u64) << 32 | (2 << 24) as u64;
//     std::env::set_var("SteamGameId", game_id.to_string());
//     std::env::set_var("SteamOverlayGameId", game_id.to_string());
//     // TODO: is this needed? decode the values
//     // std::env::set_var("EnableConfiguratorSupport", "4111");
//     std::env::set_var("SteamPath", steam_path.to_string_lossy().to_string());

//     // TODO: is this always the same, and always existing?
//     let gamepad_info_path = steam_path
//         .clone()
//         .join("config")
//         .join("virtualgamepadinfo.txt");
//     if !gamepad_info_path.exists() {
//         warn!(
//             "Steam virtualgamepadinfo.txt not found at expected path: {}",
//             gamepad_info_path.display()
//         );
//         return Err(anyhow::anyhow!("Steam virtualgamepadinfo.txt not found"));
//     }
//     // Is needed for steamHandles to be created
//     std::env::set_var(
//         "SteamVirtualGamepadInfo",
//         gamepad_info_path.to_string_lossy().to_string(),
//     );
// }
// Ok(())
