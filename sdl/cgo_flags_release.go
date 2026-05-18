//go:build release

package sdl

/*
#cgo LDFLAGS: -L${SRCDIR}/../deps/SDL/build/Release -lSDL3
*/
import "C"
