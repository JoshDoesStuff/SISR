package config

import (
	"context"
	"net/http"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/config"
	"github.com/danielgtaylor/huma/v2"
)

type ConfigResponse struct {
	Body Config
}

type Config struct {
	AutoUpdate             *config.AutoUpdate             `json:"autoUpdate"`
	RunMisc                *config.RunMisc                `json:"runMisc"`
	ControllerEmulation    *config.ControllerEmulation    `json:"controllerEmulation"`
	KeyboardMouseEmulation *config.KeyboardMouseEmulation `json:"keyboardMouseEmulation"`
	Viiper                 *config.Viiper                 `json:"viiper"`
	Window                 *config.Window                 `json:"window"`
	Steam                  *config.Steam                  `json:"steam"`
}

func Register(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method: http.MethodGet,
		Path:   "/api/v1/config",
	}, getConfig(c))
}

func getConfig(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*ConfigResponse, error) {
	return func(ctx context.Context, req *struct{}) (*ConfigResponse, error) {

		c.Config.Lock()
		defer c.Config.Unlock()

		return &ConfigResponse{
			Body: Config{
				AutoUpdate:             c.Config.AutoUpdate,
				RunMisc:                c.Config.RunMisc,
				ControllerEmulation:    c.Config.ControllerEmulation,
				KeyboardMouseEmulation: c.Config.KeyboardMouseEmulation,
				Viiper: &config.Viiper{
					Address:  c.Config.Viiper.Address, // nolint
					Password: "REDACTED",
				},
				Window: c.Config.Window,
				Steam:  c.Config.Steam,
			},
		}, nil
	}
}
