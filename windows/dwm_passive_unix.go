//go:build !windows

package windows

import "github.com/Alia5/SISR/sdl"

func SetDWMPassiveUpdateMode(window *sdl.Window) error {
	return nil
}
