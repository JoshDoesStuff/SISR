package sdl

/*
#cgo CFLAGS: -I${SRCDIR}/../deps/SDL/include

#include <stdlib.h>

#include <SDL3/SDL_video.h>
#include <SDL3/SDL_render.h>
*/
import "C"

// Renderer is a 2D rendering context.
type Renderer interface {
	// Destroy releases the rendering context.
	Destroy()
	// SetRenderDrawColor sets the color used for drawing operations.
	//
	// Set the color for drawing or filling rectangles, lines, and points, and for
	// SDL_RenderClear().
	SetRenderDrawColor(r, g, b, a uint8) error
	// RenderClear clears the current rendering target with the drawing color.
	RenderClear() error
	// RenderPresent updates the screen with any rendering performed since the previous call.
	RenderPresent() error
}

type cRenderer struct {
	cRenderer *C.SDL_Renderer
}

func (r *cRenderer) Destroy() {
	C.SDL_DestroyRenderer(r.cRenderer)
	r.cRenderer = nil
}

// SetRenderDrawColor sets the color used for drawing operations.
func (r *cRenderer) SetRenderDrawColor(rC, gC, bC, aC uint8) error {
	if !C.SDL_SetRenderDrawColor(r.cRenderer, C.Uint8(rC), C.Uint8(gC), C.Uint8(bC), C.Uint8(aC)) {
		return GetError()
	}
	return nil
}

// RenderClear clears the current rendering target with the drawing color.
func (r *cRenderer) RenderClear() error {
	if !C.SDL_RenderClear(r.cRenderer) {
		return GetError()
	}
	return nil
}

// RenderPresent updates the screen with any rendering performed since the previous call.
func (r *cRenderer) RenderPresent() error {
	if !C.SDL_RenderPresent(r.cRenderer) {
		return GetError()
	}
	return nil
}
