package ui

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/sdl/extras"
	"github.com/Alia5/SISR/webview"
	"github.com/danielgtaylor/huma/v2"
)

type ShowHideUIRequest struct {
	Body ShowHideUI
}

type ShowHideUI struct {
	Show bool `json:"show"`
}

func Register(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/ui",
	}, showHideUI(c))
}

func showHideUI(c *cmd.SISRContext) func(ctx context.Context, req *ShowHideUIRequest) (*struct{}, error) {
	return func(ctx context.Context, req *ShowHideUIRequest) (*struct{}, error) {

		c.Config.Lock()
		fullscreen := c.Config.Fullscreen
		c.Config.Unlock()

		_, err := cmd.ScheduleWindowDispatch(ctx, c.WindowDispatcher, func(w *sdl.Window, wv webview.WebView) bool {
			if req.Body.Show {
				w.ShowWindow()
				_ = c.WindowDispatcher.Schedule(func(w *sdl.Window, wv webview.WebView) any {
					wv.SetVisible(true)
					return nil
				})
				err := extras.SetCursorHitTest(w, true)
				if err != nil {
					slog.Error("Failed setting window cursor hittest", "error", err)
				}
				return true
			} else {
				if !fullscreen {
					w.HideWindow()
				}
				wv.SetVisible(false)
				err := extras.SetCursorHitTest(w, false)
				if err != nil {
					slog.Error("Failed setting window cursor hittest", "error", err)
				}
				return false
			}
		})
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
