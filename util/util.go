//go:build !windows

package util

func IsRunFromGUI() bool {
	// On non-Windows, always return false.
	// We only use this to get the bost of both worlds, using the cli, and "pretending" to be linked as native GUI APP
	return false
}

func HideConsoleWindow() {
	// No-op on non-Windows
}
