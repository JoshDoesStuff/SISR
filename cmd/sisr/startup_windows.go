//go:build windows

package main

import "github.com/Alia5/SISR/windows"

func init() {
	if windows.IsRunFromGUI() {
		windows.HideConsoleWindow()
	}
}
