package helper

import (
	"os"
	"path/filepath"
)

func GetDataDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	if configDir == "" {
		return "", os.ErrNotExist
	}
	return filepath.Join(configDir, "SISR", "data"), nil
}
