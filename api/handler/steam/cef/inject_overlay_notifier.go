package cef

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Alia5/SISR/cefpayloads"
	"github.com/Alia5/SISR/cmd"
	"github.com/danielgtaylor/huma/v2"
)

func registerInjectOverlayNotifier(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/steam/cef/inject-overlay-notifier",
		Tags:   []string{"steam", "cef", "inject"},
	}, injectOverlayNotifier(c))
}

func injectOverlayNotifier(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*struct{}, error) {
	return func(ctx context.Context, req *struct{}) (*struct{}, error) {
		c.Config.Lock()
		steamCfg := *c.Config.Steam
		c.Config.Unlock()

		slog.Debug("APIAddr", "value", *c.APIAddr)
		_, err := cefpayloads.NewOverlayCallback(&steamCfg).Execute(ctx, &cefpayloads.OverlayCallbackArgs{
			SISRAPIURL: *c.APIAddr,
		})
		if err != nil {
			slog.Error("failed to inject overlay notifier", "error", err)
		}

		return nil, err
	}
}
