package sdl

/*
#cgo CFLAGS: -I${SRCDIR}/../deps/SDL/include

#include <stdlib.h>

#include <SDL3/SDL_gamepad.h>
#include <SDL3/SDL_properties.h>
*/
import "C"

import "unsafe"

// GetBooleanProperty gets a boolean property from an SDL properties group.
func GetBooleanProperty(props uintptr, name string, defaultValue bool) bool {
	if props == 0 || name == "" {
		return defaultValue
	}

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	return bool(C.SDL_GetBooleanProperty(C.SDL_PropertiesID(props), cName, C.bool(defaultValue)))
}

// PropGamepadCapMonoLEDBoolean returns the gamepad mono LED capability property key.
func PropGamepadCapMonoLEDBoolean() string {
	return C.SDL_PROP_GAMEPAD_CAP_MONO_LED_BOOLEAN
}

// PropGamepadCapRGBLEDBoolean returns the gamepad RGB LED capability property key.
func PropGamepadCapRGBLEDBoolean() string {
	return C.SDL_PROP_GAMEPAD_CAP_RGB_LED_BOOLEAN
}

// PropGamepadCapPlayerLEDBoolean returns the gamepad player LED capability property key.
func PropGamepadCapPlayerLEDBoolean() string {
	return C.SDL_PROP_GAMEPAD_CAP_PLAYER_LED_BOOLEAN
}
