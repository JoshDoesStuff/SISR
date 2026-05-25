package cefpayloads

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"
	"text/template"

	"github.com/Alia5/SISR/config"
	steamcef "github.com/Alia5/SISR/steam/cefinject"
)

var cleanupTabs = []string{
	"SharedJSContext",
}

type SISRCleanupExecutor interface {
	steamcef.Executor[*struct{}, *struct{}]
}

//go:embed dist/cleanup.js.tmpl
var cleanupTmpl string

func NewSISRCleanup(cfg *config.Steam) SISRCleanupExecutor {
	return steamcef.NewExecutor[*struct{}, *struct{}](
		cfg,
		template.Must(template.New("cleanup").
			Delims("<<%", "%>>").Parse(cleanupTmpl)),
	)
}

func SISRCleanup(ctx context.Context, cfg *config.Steam) error {
	executor := NewSISRCleanup(cfg)
	var errs []error
	for _, tab := range cleanupTabs {
		_, err := executor.ExecuteInTab(ctx, tab, &struct{}{})
		if err != nil {
			slog.Warn("cleanup failed for tab", "tab", tab, "err", err)
			errs = append(errs, err)
		}
	}
	if len(errs) == len(cleanupTabs) {
		return errors.Join(errs...)
	}
	return nil
}
