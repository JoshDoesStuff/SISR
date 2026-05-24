//go:build !windows

package helper

/*
#cgo linux LDFLAGS: -ldl
#cgo freebsd LDFLAGS: -ldl
#cgo openbsd LDFLAGS: -ldl
#cgo netbsd LDFLAGS: -ldl
#include <dlfcn.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

func LoadLib(path string) error {
	cPath := C.CString(path)
	defer C.free(unsafe.Pointer(cPath))

	handle := C.dlopen(cPath, C.RTLD_NOW|C.RTLD_GLOBAL)
	if handle == nil {
		err := C.dlerror()
		if err == nil {
			return errors.New("dlopen failed")
		}

		return errors.New(C.GoString(err))
	}

	return nil
}
