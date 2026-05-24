package steam

import (
	"errors"
	"fmt"
	"os"
	"path"
)

func ShortcutsPath(steamDir string, userID uint32) (string, error) {
	if steamDir == "" {
		var err error
		steamDir, err = steamPath()
		if err != nil {
			return "", err
		}
	}
	if userID == 0 {
		var err error
		userID, err = ActiveUserID()
		if err != nil {
			return "", err
		}
	}

	shortcutsPath := path.Join(
		steamDir,
		"userdata",
		fmt.Sprintf("%d", userID),
		"config",
		"shortcuts.vdf",
	)
	if _, err := os.Stat(shortcutsPath); err != nil {
		return "", errors.New("shortcuts.vdf does not exist")
	}

	return shortcutsPath, nil
}
