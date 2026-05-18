//go:build !release

package sdl

/*
#cgo LDFLAGS: -L${SRCDIR}/../deps/SDL/build/Debug -lSDL3
*/
import "C"
