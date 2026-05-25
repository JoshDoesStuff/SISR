package api

import (
	"github.com/Alia5/SISR/api/handler/config"
	"github.com/Alia5/SISR/api/handler/devices"
	"github.com/Alia5/SISR/api/handler/initiallaunch"
	"github.com/Alia5/SISR/api/handler/quit"
	"github.com/Alia5/SISR/api/handler/steam"
	"github.com/Alia5/SISR/api/handler/ui"
	"github.com/Alia5/SISR/api/handler/version"
	"github.com/Alia5/SISR/api/handler/viiper"
	"github.com/Alia5/SISR/cmd"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterAPI(hAPI huma.API, c *cmd.SISRContext) {
	quit.Register(hAPI, c)
	steam.Register(hAPI, c)
	devices.Register(hAPI, c)
	viiper.Register(hAPI, c)
	initiallaunch.Register(hAPI, c)
	version.Register(hAPI, c)
	ui.Register(hAPI, c)
	config.Register(hAPI, c)
}
