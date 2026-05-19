package polyhook

/*
#cgo CXXFLAGS: -std=c++20
#cgo CXXFLAGS: -I${SRCDIR}/../deps/PolyHook2
#cgo CXXFLAGS: -I${SRCDIR}/../deps/PolyHook2/asmjit/src
#cgo CXXFLAGS: -I${SRCDIR}/../deps/PolyHook2/asmtk/src
#cgo CXXFLAGS: -I${SRCDIR}/../deps/PolyHook2/zydis/include
#cgo CXXFLAGS: -I${SRCDIR}/../deps/PolyHook2/zydis/dependencies/zycore/include
#cgo LDFLAGS: -L${SRCDIR}/../deps/PolyHook2/build/Debug -lPolyHook_2
#cgo LDFLAGS: -L${SRCDIR}/../deps/PolyHook2/build/zydis/Debug -lZydis
#cgo LDFLAGS: -L${SRCDIR}/../deps/PolyHook2/build/zydis/dependencies/zycore/Debug -lZycore
#cgo LDFLAGS: -L${SRCDIR}/../deps/PolyHook2/build/asmtk/Debug -lasmjit -lasmtk

#include "polyhook.h"
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Detour represents an x64 function detour hook backed by PolyHook2.
type Detour struct {
	d          *C.PLH_Detour
	Trampoline uintptr // populated after a successful Hook() call
}

// NewDetour creates a new x64 detour.
//
// fnAddress is the address of the target function to hook.
// fnCallback is the address of the replacement function.
//
// The returned Detour must be closed with Close() when no longer needed.
// Call Hook() to actually install the detour.
func NewDetour(fnAddress, fnCallback uintptr) (*Detour, error) {
	d := C.PLH_x64Detour_new(C.uint64_t(fnAddress), C.uint64_t(fnCallback))
	if d == nil {
		return nil, errors.New("polyhook: failed to allocate x64Detour")
	}
	return &Detour{d: d}, nil
}

// Hook installs the detour in memory.
// After a successful call, Trampoline holds the address of the original function.
func (d *Detour) Hook() error {
	if C.PLH_x64Detour_hook(d.d) == 0 {
		return detourError(d.d)
	}
	d.Trampoline = uintptr(C.PLH_x64Detour_trampoline(d.d))
	return nil
}

// Unhook removes the detour and restores the original function bytes.
func (d *Detour) Unhook() error {
	if C.PLH_x64Detour_unhook(d.d) == 0 {
		return detourError(d.d)
	}
	return nil
}

// Close frees all PolyHook2 resources associated with this detour.
// Unhook() should be called before Close() if the hook is currently active.
func (d *Detour) Close() {
	if d.d != nil {
		C.PLH_x64Detour_free(d.d)
		d.d = nil
	}
}

// detourError reads the last error string from PolyHook2.
func detourError(d *C.PLH_Detour) error {
	msg := C.GoString(C.PLH_x64Detour_last_error(d))
	if msg == "" {
		return errors.New("polyhook: operation failed (no details)")
	}
	return errors.New("polyhook: " + msg)
}

// UnsafePointerToUintptr converts an unsafe.Pointer to a uintptr for use as
// a hook address.  This is a convenience helper; the caller is responsible for
// ensuring the pointer remains valid.
func UnsafePointerToUintptr(p unsafe.Pointer) uintptr {
	return uintptr(p)
}
