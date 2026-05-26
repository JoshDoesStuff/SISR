package steam

import (
	"net/http"

	"github.com/Alia5/SISR/api/handler/steam/cef"
	"github.com/Alia5/SISR/cmd"
	"github.com/danielgtaylor/huma/v2"
)

func Register(a huma.API, c *cmd.SISRContext) {
	huma.Register(a, huma.Operation{
		Method: http.MethodGet,
		Path:   "/api/v1/steam/status",
	}, status(c))

	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/steam/enable-cef-remote-debugging",
	}, enableCEFRemoteDebugging(c))

	huma.Register(a, huma.Operation{
		Method: http.MethodPost,
		Path:   "/api/v1/steam/restart",
	}, restartSteam(c))

	cef.Register(a, c)

}
