package steam

import (
	"context"
	"log/slog"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/steam"
)

func enableCEFRemoteDebugging(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*struct{}, error) {
	return func(ctx context.Context, req *struct{}) (*struct{}, error) {

		c.Config.Lock()
		steamCfg := *c.Config.Steam
		c.Config.Unlock()

		_, err := steam.EnableCefRemoteDebug(&steamCfg)
		if err != nil {
			slog.Error("Failed to enable CEF remote debugging", "error", err)
		}

		return &struct{}{}, err
	}
}
