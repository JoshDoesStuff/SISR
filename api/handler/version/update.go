package version

import (
	"context"
	"net/http"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/meta"
	"github.com/danielgtaylor/huma/v2"
)

type VersionInfoResponse struct {
	Body VersionInfo
}

type VersionInfo struct {
	Version         string `json:"version"`
	UpdateAvailable bool   `json:"update_available"`
	NewVersion      string `json:"new_version,omitempty"`
}

func Register(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method: http.MethodGet,
		Path:   "/api/v1/version/info",
	}, updateAvailable(c))
	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/version/update_ceck",
	}, checkForUpdate(c))
}

func updateAvailable(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*VersionInfoResponse, error) {
	return func(ctx context.Context, req *struct{}) (*VersionInfoResponse, error) {
		versionInfo := c.UpdateChecker.GetVersionInfo()
		return &VersionInfoResponse{
			Body: VersionInfo{
				Version:         versionInfo.Version,
				UpdateAvailable: versionInfo.UpdateAvailable,
				NewVersion:      versionInfo.NewVersion,
			},
		}, nil
	}
}

func checkForUpdate(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*VersionInfoResponse, error) {
	return func(ctx context.Context, req *struct{}) (*VersionInfoResponse, error) {
		updateInfo, err := c.UpdateChecker.CheckForUpdate(ctx)
		if err != nil {
			return nil, err
		}
		return &VersionInfoResponse{
			Body: VersionInfo{
				Version:         meta.Version,
				UpdateAvailable: updateInfo.UpdateAvailable,
				NewVersion:      updateInfo.NewVersion,
			},
		}, nil
	}
}
