package gamepadUpdated

import (
	"encoding"
	"math"

	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/VIIPER/device/xbox360"
)

func toXbox360State(gp *sdl.Gamepad, state *encoding.BinaryMarshaler) {
	s, ok := (*state).(*xbox360.InputState)
	if !ok || s == nil {
		s = &xbox360.InputState{}
		*state = s
	}
	s.Buttons = 0
	if gp.GetButton(sdl.GamepadButtonSouth) {
		s.Buttons |= xbox360.ButtonA
	}
	if gp.GetButton(sdl.GamepadButtonEast) {
		s.Buttons |= xbox360.ButtonB
	}
	if gp.GetButton(sdl.GamepadButtonWest) {
		s.Buttons |= xbox360.ButtonX
	}
	if gp.GetButton(sdl.GamepadButtonNorth) {
		s.Buttons |= xbox360.ButtonY
	}
	if gp.GetButton(sdl.GamepadButtonLeftShoulder) {
		s.Buttons |= xbox360.ButtonLShoulder
	}
	if gp.GetButton(sdl.GamepadButtonRightShoulder) {
		s.Buttons |= xbox360.ButtonRShoulder
	}
	if gp.GetButton(sdl.GamepadButtonLeftStick) {
		s.Buttons |= xbox360.ButtonLThumb
	}
	if gp.GetButton(sdl.GamepadButtonRightStick) {
		s.Buttons |= xbox360.ButtonRThumb
	}
	if gp.GetButton(sdl.GamepadButtonStart) {
		s.Buttons |= xbox360.ButtonStart
	}
	if gp.GetButton(sdl.GamepadButtonBack) {
		s.Buttons |= xbox360.ButtonBack
	}
	if gp.GetButton(sdl.GamepadButtonGuide) {
		s.Buttons |= xbox360.ButtonGuide
	}
	if gp.GetButton(sdl.GamepadButtonDpadUp) {
		s.Buttons |= xbox360.ButtonDPadUp
	}
	if gp.GetButton(sdl.GamepadButtonDpadDown) {
		s.Buttons |= xbox360.ButtonDPadDown
	}
	if gp.GetButton(sdl.GamepadButtonDpadLeft) {
		s.Buttons |= xbox360.ButtonDPadLeft
	}
	if gp.GetButton(sdl.GamepadButtonDpadRight) {
		s.Buttons |= xbox360.ButtonDPadRight
	}

	lt := gp.GetAxis(sdl.GamepadAxisLeftTrigger)
	rt := gp.GetAxis(sdl.GamepadAxisRightTrigger)

	s.LT = uint8(max(0, min(math.MaxUint8, max(0, int32(lt))*math.MaxUint8/math.MaxInt16)))
	s.RT = uint8(max(0, min(math.MaxUint8, max(0, int32(rt))*math.MaxUint8/math.MaxInt16)))

	// Invert Y axes to match XInput convention
	// XInput: Negative values signify down or to the left. Positive values signify up or to the right.
	//         https://learn.microsoft.com/en-us/windows/win32/api/xinput/ns-xinput-xinput_gamepad
	// SDL: For thumbsticks, the state is a value ranging from -32768 (up/left) to 32767 (down/right).
	//      https://wiki.libsdl.org/SDL3/SDL_GetGamepadAxis
	s.LX = gp.GetAxis(sdl.GamepadAxisLeftX)
	s.LY = int16(max(math.MinInt16, min(math.MaxInt16, -int32(gp.GetAxis(sdl.GamepadAxisLeftY)))))
	s.RX = gp.GetAxis(sdl.GamepadAxisRightX)
	s.RY = int16(max(math.MinInt16, min(math.MaxInt16, -int32(gp.GetAxis(sdl.GamepadAxisRightY)))))

}
