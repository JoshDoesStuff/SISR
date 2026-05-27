package cef

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Alia5/SISR/cefpayloads"
	"github.com/Alia5/SISR/cmd"
	"github.com/danielgtaylor/huma/v2"
)

type OpenConfiguratorRequest struct {
	Body cefpayloads.OpenConfiguratorArgs
}

func registerOpenConfigurator(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/steam/cef/open-configurator",
		Tags:   []string{"steam", "cef", "layout-config"},
	}, openConfigurator(c))
}

func openConfigurator(c *cmd.SISRContext) func(ctx context.Context, req *OpenConfiguratorRequest) (*struct{}, error) {
	return func(ctx context.Context, req *OpenConfiguratorRequest) (*struct{}, error) {
		c.Config.Lock()
		steamCfg := *c.Config.Steam
		c.Config.Unlock()

		_, err := cefpayloads.NewOpenConfigurator(&steamCfg).Execute(ctx, &req.Body)
		if err != nil {
			slog.Error("failed to open configurator", "error", err)
		}

		return nil, err
	}
}
