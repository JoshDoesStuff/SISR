package sdl

// WindowPointerProperty selects a platform window pointer property.
type WindowPointerProperty int

// Known platform window pointer properties.
const (
	WindowPointerPropertyWin32HWND WindowPointerProperty = iota
	WindowPointerPropertyWin32HDC
	WindowPointerPropertyWin32Instance
	WindowPointerPropertyWaylandDisplay
	WindowPointerPropertyWaylandSurface
	WindowPointerPropertyX11Display
)
