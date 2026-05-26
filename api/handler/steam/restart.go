package steam

import (
	"context"
	"errors"
	"time"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/helper"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/steam"
	"github.com/Alia5/SISR/webview"
)

func restartSteam(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*struct{}, error) {
	return func(ctx context.Context, req *struct{}) (*struct{}, error) {

		_ = c.WindowDispatcher.Schedule(func(w *sdl.Window, wv webview.WebView) any {
			_ = w.SetWindowAlwaysOnTop(false)
			return nil
		})

		steamRunning := steam.ClientRunning()
		if steamRunning {
			err := helper.OpenURL("steam://exit")
			if err != nil {
				return nil, err
			}
		}

		steamShutdown := false
		for range 30 {
			steamRunning := steam.ClientRunning()
			if !steamRunning {
				steamShutdown = true
				break
			}
			time.Sleep(1 * time.Second)
		}
		if !steamShutdown {
			return nil, errors.New("steam did not shut down within the expected time")
		}
		err := helper.OpenURL("steam://open/main")
		if err != nil {
			return nil, err
		}
		steamStart := false
		for range 30 {
			steamRunning := steam.ClientRunning()
			if steamRunning {
				steamStart = true
				break
			}
			time.Sleep(1 * time.Second)
		}
		if !steamStart {
			return nil, errors.New("steam did not start within the expected time")
		}
		// hack, give Steam a bit more time....
		time.Sleep(15 * time.Second)
		userLoggedIn := false
		for range 30 {
			userID, err := steam.ActiveUserID()
			if err == nil && userID != 0 {
				userLoggedIn = true
				break
			}
			time.Sleep(1 * time.Second)
		}
		if !userLoggedIn {
			return nil, errors.New("user did not log in within the expected time")
		}
		// hack, give Steam even more time....
		time.Sleep(5 * time.Second)

		return nil, nil
	}
}
