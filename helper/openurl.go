package helper

import (
	"log/slog"
	"os/exec"
	"runtime"
	"strings"
)

func OpenURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "rundll32.exe"
		args = []string{"url.dll,FileProtocolHandler", strings.ReplaceAll(url, "&", "^&")}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}

	slog.Debug("Opening URL", "url", url, "cmd", cmd, "args", args)
	e := exec.Command(cmd, args...)
	err := e.Start()
	if err != nil {
		return err
	}
	err = e.Wait()
	if err != nil {
		return err
	}

	return nil
}
