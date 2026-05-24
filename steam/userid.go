//go:build !windows

package steam

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func ActiveUserID() (uint32, error) {
	sp, err := steamPath()
	if err != nil {
		return 0, err
	}

	registryVDF := filepath.Join(filepath.Dir(sp), "registry.vdf")
	registryVDFInfo, statErr := os.Stat(registryVDF)
	if statErr == nil && !registryVDFInfo.IsDir() {
		content, readErr := os.ReadFile(registryVDF)
		if readErr == nil {
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				if !strings.HasPrefix(trimmed, `"ActiveUser"`) {
					continue
				}

				parts := strings.Split(trimmed, "\"")
				if len(parts) < 4 {
					continue
				}

				userID, parseErr := strconv.ParseUint(parts[3], 10, 32)
				if parseErr != nil || userID == 0 {
					continue
				}

				slog.Debug("Found active Steam user ID from registry.vdf", "userID", userID)
				return uint32(userID), nil
			}
		}
	}

	userdataPath := filepath.Join(sp, "userdata")
	userdataInfo, statErr := os.Stat(userdataPath)
	if statErr == nil && userdataInfo.IsDir() {
		entries, readDirErr := os.ReadDir(userdataPath)
		if readDirErr == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}

				userID, parseErr := strconv.ParseUint(entry.Name(), 10, 32)
				if parseErr != nil || userID == 0 {
					continue
				}

				slog.Debug("Found possibly active Steam user ID from userdata directory", "userID", userID)
				return uint32(userID), nil
			}
		}
	}

	return 0, fmt.Errorf("active Steam user ID not found")
}
