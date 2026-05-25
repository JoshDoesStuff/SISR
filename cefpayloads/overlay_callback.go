package cefpayloads

import (
	_ "embed"
	"text/template"

	"github.com/Alia5/SISR/config"
	steamcef "github.com/Alia5/SISR/steam/cefinject"
)

type OverlayCallbackArgs struct {
	SISRAPIURL string `json:"SISRAPIURL" required:"true"`
}

type OverlayCallbackExecutor interface {
	steamcef.Executor[*OverlayCallbackArgs, *struct{}]
}

//go:embed dist/overlay_callback.js.tmpl
var overlayCallbackJSTmpl string
var overlayCallbackJS = template.Must(
	template.New("overlayCallback").
		Delims("<<%", "%>>").Parse(overlayCallbackJSTmpl),
)

func NewOverlayCallback(cfg *config.Steam) OverlayCallbackExecutor {
	return steamcef.NewExecutor[*OverlayCallbackArgs, *struct{}](cfg, overlayCallbackJS)
}
