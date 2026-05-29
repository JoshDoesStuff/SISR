package helper

import (
	"os"
	"path/filepath"
)

func GetOwnExecutablePath() (string, error) {
	runningPath := os.Getenv("APPIMAGE")
	if runningPath == "" {
		var err error
		runningPath, err = os.Executable()
		if err != nil {
			return "", err
		}
	}

	if resolvedPath, err := filepath.EvalSymlinks(runningPath); err == nil {
		runningPath = resolvedPath
	}

	if absolutePath, err := filepath.Abs(runningPath); err == nil {
		runningPath = absolutePath
	}

	return runningPath, nil
}

func GetOwnExecutableDir() (string, error) {
	exePath, err := GetOwnExecutablePath()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}
