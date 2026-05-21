//go:build windows

package config

type PlatformOpts struct {
	Console bool `help:"Show console window" default:"false" env:"SISR_CONSOLE"`
}
