package api

import (
	"github.com/Alia5/SISR/api/handler"
	"github.com/Alia5/SISR/cmd"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterAPI(hAPI huma.API, c *cmd.SISRContext) {
	handler.RegisterQuit(hAPI, c)
	handler.RegisterSteamStatus(hAPI, c)
}
