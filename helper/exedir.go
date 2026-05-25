package helper

import (
	"log/slog"
	"os"
	"path/filepath"
)

func GetOwnExecutableDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	ownExeDir, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		slog.Error("Failed to evaluate symlinks for own executable path", "error", err)
		return "", err
	}
	return filepath.Dir(ownExeDir), nil
}
