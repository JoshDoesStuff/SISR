package cef

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

func registerOverlayStateChanged(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method:      http.MethodPost,
		Path:        "/api/v1/steam/cef/overlay-state-changed",
		Description: "Callback for CEF Injectee when Steam overlay is opened/closed",
		Tags:        []string{"steam", "cef", "callback"},
	}, overlayStateChanged(c))
}

type OverlayChangedRequest struct {
	Body OverlayChanged
}

type OverlayChanged struct {
	Open bool `json:"open"`
}

func overlayStateChanged(c *cmd.SISRContext) func(ctx context.Context, req *OverlayChangedRequest) (*struct{}, error) {

	previousVisibility := false
	previousHitTest := false

	return func(ctx context.Context, req *OverlayChangedRequest) (*struct{}, error) {
		c.Config.Lock()
		kbmEnabled := c.Config.KeyboardMouseEmulation
		c.Config.Unlock()

		err, dispatchErr := cmd.ScheduleWindowDispatch(
			ctx,
			c.WindowDispatcher,
			func(w *sdl.Window, wv webview.WebView) error {
				if req.Body.Open {
					previousVisibility = wv.Visible()
					previousHitTest = extras.GetWindowHitTest(w)
				}

				slog.Debug("Overlay state changed",
					"previousWebViewVisibility", previousVisibility,
					"previousCursorHitTest", previousHitTest,
					"newOverlayOpenState", req.Body.Open,
				)

				if req.Body.Open {
					wv.SetVisible(false)
					err := extras.SetCursorHitTest(w, true)
					if err != nil {
						slog.Error("Failed setting window cursor hittest", "error", err)
					}
				} else {
					if previousVisibility {
						wv.SetVisible(true)
					}
					targetHitTest := previousHitTest
					if kbmEnabled {
						targetHitTest = true
					}
					err := extras.SetCursorHitTest(w, targetHitTest)
					if err != nil {
						slog.Error("Failed setting window cursor hittest", "error", err)
					}
				}
				return nil
			})
		if dispatchErr != nil {
			return nil, dispatchErr
		}
		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
