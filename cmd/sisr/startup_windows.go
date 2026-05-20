//go:build windows

package main

import (
	"github.com/Alia5/SISR/util"
)

func init() {
	if util.IsRunFromGUI() {
		util.HideConsoleWindow()
	}
}
