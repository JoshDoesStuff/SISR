//go:build windows

package main

import (
	"github.com/Alia5/SISR/config"
	"github.com/Alia5/SISR/windows"
)

func applyPlatformStartup(cfg config.Global) {
	if cfg.Console {
		return
	}

	if windows.IsRunFromGUI() {
		windows.HideConsoleWindow()
	}
}
