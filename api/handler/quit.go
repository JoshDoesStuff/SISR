package handler

import (
	"context"
	"net/http"

	"github.com/Alia5/SISR/cmd"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterQuit(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/shutdown",
	}, quitHanleFunc(c))
}

func quitHanleFunc(c *cmd.SISRContext) func(ctx context.Context, req *struct{}) (*struct{}, error) {
	return func(ctx context.Context, req *struct{}) (*struct{}, error) {
		defer func() {
			c.QuitFn()
		}()
		return nil, nil
	}
}
