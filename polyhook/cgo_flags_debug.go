//go:build !release

package polyhook

/*
#cgo LDFLAGS: -L${SRCDIR}/../deps/PolyHook2/build/Debug -lPolyHook_2
#cgo LDFLAGS: -L${SRCDIR}/../deps/PolyHook2/build/zydis/Debug -lZydis
#cgo LDFLAGS: -L${SRCDIR}/../deps/PolyHook2/build/zydis/dependencies/zycore/Debug -lZycore
#cgo LDFLAGS: -L${SRCDIR}/../deps/PolyHook2/build/asmtk/Debug -lasmjit -lasmtk
*/
import "C"
