package cli

import (
	"github.com/Alia5/SISR/cli/cmd"
	"github.com/Alia5/SISR/config"
)

type CLI struct {
	Config config.Global `embed:""`
	SISR   cmd.SISR      `cmd:"" help:"Run SISR (Default)" default:"1" name:" "`
}
