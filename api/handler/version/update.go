package version

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/helper"
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
	UpdateDismissed bool   `json:"update_dismissed"`
}

func Register(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method: http.MethodGet,
		Path:   "/api/v1/version/info",
	}, updateAvailable(c))

	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/version/update-check",
	}, checkForUpdate(c))

	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/version/skip-update",
	}, skipUpdate(c))

	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/version/update-remind-later",
	}, remindLater(c))

	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/version/update-view-on-github",
	}, viewOnGitHub(c))
	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/version/install-update",
	}, installUpdate(c))
}

func updateAvailable(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*VersionInfoResponse, error) {
	return func(ctx context.Context, req *struct{}) (*VersionInfoResponse, error) {
		versionInfo := c.UpdateChecker.GetVersionInfo()
		return &VersionInfoResponse{
			Body: VersionInfo{
				Version:         versionInfo.Version,
				UpdateAvailable: versionInfo.UpdateAvailable,
				NewVersion:      versionInfo.NewVersion,
				UpdateDismissed: versionInfo.UpdateDismissed,
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
				UpdateDismissed: updateInfo.UpdateDismissed,
			},
		}, nil
	}
}

func skipUpdate(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*struct{}, error) {
	return func(ctx context.Context, req *struct{}) (*struct{}, error) {
		err := c.UpdateChecker.SkipUpdate()
		return nil, err
	}
}

func remindLater(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*struct{}, error) {
	return func(ctx context.Context, req *struct{}) (*struct{}, error) {
		err := c.UpdateChecker.RemindLater()
		return nil, err
	}
}

func viewOnGitHub(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*struct{}, error) {
	return func(ctx context.Context, req *struct{}) (*struct{}, error) {
		err := helper.OpenURL(fmt.Sprintf("https://github.com/Alia5/SISR/releases/tag/%s", c.UpdateChecker.GetVersionInfo().NewVersion))
		return nil, err
	}
}

func installUpdate(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*struct{}, error) {
	return func(ctx context.Context, req *struct{}) (*struct{}, error) {
		// TODO
		return nil, errors.New("not implemented")
	}
}
