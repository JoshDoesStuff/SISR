package restart

import (
	"context"
	"net/http"
	"os"
	"os/exec"

	"github.com/Alia5/SISR/cmd"
	"github.com/Alia5/SISR/helper"
	"github.com/danielgtaylor/huma/v2"
)

func Register(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/restart-sisr",
	}, restart(c))
}

func restart(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*struct{}, error) {
	return func(ctx context.Context, req *struct{}) (*struct{}, error) {

		ownExecutable, err := helper.GetOwnExecutablePath()
		if err != nil {
			return nil, err
		}
		args := os.Args[1:]
		cmd := exec.Command(ownExecutable, args...)
		err = cmd.Start()
		if err != nil {
			return nil, err
		}
		defer func() {
			c.QuitFn()
		}()

		return nil, nil
	}
}
