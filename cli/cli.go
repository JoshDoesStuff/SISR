package cli

import (
	"github.com/Alia5/SISR/cli/cmd/sisr"
	"github.com/Alia5/SISR/config"
)

type CLI struct {
	Config config.Global `embed:""`
	SISR   sisr.SISR     `cmd:"" help:"Run SISR (Default)" default:"1" name:"run"`
}
