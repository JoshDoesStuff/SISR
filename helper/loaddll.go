//go:build windows

package helper

import "syscall"

func LoadLib(path string) error {
	return syscall.NewLazyDLL(path).Load()
}
