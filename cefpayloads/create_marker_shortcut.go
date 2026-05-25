package cefpayloads

import (
	_ "embed"
	"text/template"

	"github.com/Alia5/SISR/config"
	steamcef "github.com/Alia5/SISR/steam/cefinject"
)

type CreateMarkerShortcutArgs struct {
	SISRPath string `json:"SISRPath" required:"true"`
}

type CreateMarkerShortcutExecutor interface {
	steamcef.Executor[*CreateMarkerShortcutArgs, *struct{}]
}

//go:embed dist/create_marker_shortcut.js.tmpl
var createMarkerShortcutJSTmpl string
var createMarkerShortcutJS = template.Must(
	template.New("createMarkerShortcut").
		Delims("<<%", "%>>").Parse(createMarkerShortcutJSTmpl),
)

func NewCreateMarkerShortcut(cfg *config.Steam) CreateMarkerShortcutExecutor {
	return steamcef.NewExecutor[*CreateMarkerShortcutArgs, *struct{}](cfg, createMarkerShortcutJS)
}
