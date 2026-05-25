package cefpayloads

import (
	_ "embed"
	"text/template"

	"github.com/Alia5/SISR/config"
	steamcef "github.com/Alia5/SISR/steam/cefinject"
)

type OpenConfiguratorArgs struct {
	AppID uint32 `json:"appId" required:"true"`
}

type OpenConfiguratorExecutor interface {
	steamcef.Executor[*OpenConfiguratorArgs, *struct{}]
}

//go:embed dist/open_configurator.js.tmpl
var openConfiguratorJSTmpl string
var openConfiguratorJS = template.Must(
	template.New("openConfigurator").
		Delims("<<%", "%>>").Parse(openConfiguratorJSTmpl),
)

func NewOpenConfigurator(cfg *config.Steam) OpenConfiguratorExecutor {
	return steamcef.NewExecutor[*OpenConfiguratorArgs, *struct{}](cfg, openConfiguratorJS)
}
