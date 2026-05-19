//go:build release

package polyhook

/*
#cgo LDFLAGS: -L${SRCDIR}/../deps/PolyHook2/build/Release -lPolyHook_2
#cgo LDFLAGS: -L${SRCDIR}/../deps/PolyHook2/build/zydis/Release -lZydis
#cgo LDFLAGS: -L${SRCDIR}/../deps/PolyHook2/build/zydis/dependencies/zycore/Release -lZycore
#cgo LDFLAGS: -L${SRCDIR}/../deps/PolyHook2/build/asmtk/Release -lasmjit -lasmtk
*/
import "C"
