package viiperdevice

import (
	"encoding"
	"math"

	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/VIIPER/device/xbox360"
)

func toXbox360State(gp *sdl.Gamepad) encoding.BinaryMarshaler {
	state := &xbox360.InputState{}

	if gp.GetButton(sdl.GamepadButtonSouth) {
		state.Buttons |= xbox360.ButtonA
	}
	if gp.GetButton(sdl.GamepadButtonEast) {
		state.Buttons |= xbox360.ButtonB
	}
	if gp.GetButton(sdl.GamepadButtonWest) {
		state.Buttons |= xbox360.ButtonX
	}
	if gp.GetButton(sdl.GamepadButtonNorth) {
		state.Buttons |= xbox360.ButtonY
	}
	if gp.GetButton(sdl.GamepadButtonLeftShoulder) {
		state.Buttons |= xbox360.ButtonLShoulder
	}
	if gp.GetButton(sdl.GamepadButtonRightShoulder) {
		state.Buttons |= xbox360.ButtonRShoulder
	}
	if gp.GetButton(sdl.GamepadButtonLeftStick) {
		state.Buttons |= xbox360.ButtonLThumb
	}
	if gp.GetButton(sdl.GamepadButtonRightStick) {
		state.Buttons |= xbox360.ButtonRThumb
	}
	if gp.GetButton(sdl.GamepadButtonStart) {
		state.Buttons |= xbox360.ButtonStart
	}
	if gp.GetButton(sdl.GamepadButtonBack) {
		state.Buttons |= xbox360.ButtonBack
	}
	if gp.GetButton(sdl.GamepadButtonGuide) {
		state.Buttons |= xbox360.ButtonGuide
	}
	if gp.GetButton(sdl.GamepadButtonDpadUp) {
		state.Buttons |= xbox360.ButtonDPadUp
	}
	if gp.GetButton(sdl.GamepadButtonDpadDown) {
		state.Buttons |= xbox360.ButtonDPadDown
	}
	if gp.GetButton(sdl.GamepadButtonDpadLeft) {
		state.Buttons |= xbox360.ButtonDPadLeft
	}
	if gp.GetButton(sdl.GamepadButtonDpadRight) {
		state.Buttons |= xbox360.ButtonDPadRight
	}

	lt := gp.GetAxis(sdl.GamepadAxisLeftTrigger)
	rt := gp.GetAxis(sdl.GamepadAxisRightTrigger)

	state.LT = uint8(max(0, min(math.MaxUint8, max(0, int32(lt))*math.MaxUint8/math.MaxInt16)))
	state.RT = uint8(max(0, min(math.MaxUint8, max(0, int32(rt))*math.MaxUint8/math.MaxInt16)))

	// Invert Y axes to match XInput convention
	// XInput: Negative values signify down or to the left. Positive values signify up or to the right.
	//         https://learn.microsoft.com/en-us/windows/win32/api/xinput/ns-xinput-xinput_gamepad
	// SDL: For thumbsticks, the state is a value ranging from -32768 (up/left) to 32767 (down/right).
	//      https://wiki.libsdl.org/SDL3/SDL_GetGamepadAxis
	state.LX = gp.GetAxis(sdl.GamepadAxisLeftX)
	state.LY = int16(max(math.MinInt16, min(math.MaxInt16, -int32(gp.GetAxis(sdl.GamepadAxisLeftY)))))
	state.RX = gp.GetAxis(sdl.GamepadAxisRightX)
	state.RY = int16(max(math.MinInt16, min(math.MaxInt16, -int32(gp.GetAxis(sdl.GamepadAxisRightY)))))

	return state
}
