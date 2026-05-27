package gamepadupdated

import (
	"encoding"
	"math"

	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/VIIPER/device/dualshock4"
)

func toDualShock4State(gp *sdl.Gamepad, state *encoding.BinaryMarshaler) {
	s, ok := (*state).(*dualshock4.InputState)
	if !ok || s == nil {
		s = &dualshock4.InputState{}
		*state = s
	}
	s.Buttons = 0
	s.DPad = 0
	if gp.GetButton(sdl.GamepadButtonSouth) {
		s.Buttons |= dualshock4.ButtonCross
	}
	if gp.GetButton(sdl.GamepadButtonEast) {
		s.Buttons |= dualshock4.ButtonCircle
	}
	if gp.GetButton(sdl.GamepadButtonWest) {
		s.Buttons |= dualshock4.ButtonSquare
	}
	if gp.GetButton(sdl.GamepadButtonNorth) {
		s.Buttons |= dualshock4.ButtonTriangle
	}
	if gp.GetButton(sdl.GamepadButtonLeftShoulder) {
		s.Buttons |= dualshock4.ButtonL1
	}
	if gp.GetButton(sdl.GamepadButtonRightShoulder) {
		s.Buttons |= dualshock4.ButtonR1
	}
	if gp.GetButton(sdl.GamepadButtonLeftStick) {
		s.Buttons |= dualshock4.ButtonL3
	}
	if gp.GetButton(sdl.GamepadButtonRightStick) {
		s.Buttons |= dualshock4.ButtonR3
	}
	if gp.GetButton(sdl.GamepadButtonStart) {
		s.Buttons |= dualshock4.ButtonOptions
	}
	if gp.GetButton(sdl.GamepadButtonBack) {
		s.Buttons |= dualshock4.ButtonShare
	}
	if gp.GetButton(sdl.GamepadButtonGuide) {
		s.Buttons |= dualshock4.ButtonPS
	}
	if gp.GetButton(sdl.GamepadButtonDpadUp) {
		s.DPad |= dualshock4.DPadUp
	}
	if gp.GetButton(sdl.GamepadButtonDpadDown) {
		s.DPad |= dualshock4.DPadDown
	}
	if gp.GetButton(sdl.GamepadButtonDpadLeft) {
		s.DPad |= dualshock4.DPadLeft
	}
	if gp.GetButton(sdl.GamepadButtonDpadRight) {
		s.DPad |= dualshock4.DPadRight
	}

	lt := gp.GetAxis(sdl.GamepadAxisLeftTrigger)
	rt := gp.GetAxis(sdl.GamepadAxisRightTrigger)

	s.L2 = uint8(max(0, min(math.MaxUint8, max(0, int32(lt))*math.MaxUint8/math.MaxInt16)))
	s.R2 = uint8(max(0, min(math.MaxUint8, max(0, int32(rt))*math.MaxUint8/math.MaxInt16)))
	if lt > 128 {
		s.Buttons |= dualshock4.ButtonL2
	}
	if rt > 128 {
		s.Buttons |= dualshock4.ButtonR2
	}

	s.LX = int8(max(math.MinInt8, min(math.MaxInt8, int32(gp.GetAxis(sdl.GamepadAxisLeftX))*(math.MaxInt8+1)/math.MaxInt16)))
	s.LY = int8(max(math.MinInt8, min(math.MaxInt8, int32(gp.GetAxis(sdl.GamepadAxisLeftY))*(math.MaxInt8+1)/math.MaxInt16)))
	s.RX = int8(max(math.MinInt8, min(math.MaxInt8, int32(gp.GetAxis(sdl.GamepadAxisRightX))*(math.MaxInt8+1)/math.MaxInt16)))
	s.RY = int8(max(math.MinInt8, min(math.MaxInt8, int32(gp.GetAxis(sdl.GamepadAxisRightY))*(math.MaxInt8+1)/math.MaxInt16)))

}
