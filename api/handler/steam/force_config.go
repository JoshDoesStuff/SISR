package steam

import (
	"context"
	"log/slog"

	"github.com/Alia5/SISR/cmd"
)

type ForceInputConfigRequest struct {
	Body EnforceInputConfig
}

type EnforceInputConfig struct {
	Enforce bool `json:"enforce"`
}

func enforceInputConfig(c *cmd.SISRContext) func(ctx context.Context, req *ForceInputConfigRequest) (*struct{}, error) {
	return func(ctx context.Context, req *ForceInputConfigRequest) (*struct{}, error) {

		c.Config.Lock()
		defer c.Config.Unlock()

		if c.Config.NoSteam {
			return &struct{}{}, nil
		}

		if req.Body.Enforce {
			err := c.BindingEnforcer.ForceOwnAppID()
			if err != nil {
				slog.Error("failed to enforce input appID", "error", err)
				return nil, err
			}
		} else {
			err := c.BindingEnforcer.ForceInputAppID(0)
			if err != nil {
				slog.Error("failed to enforce input appID", "error", err)
				return nil, err
			}
		}

		return nil, nil
	}
}
