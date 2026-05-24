package viiperdevice

import (
	"encoding"
	"math"

	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/VIIPER/device/dualshock4"
)

func toDualShock4State(gp *sdl.Gamepad) encoding.BinaryMarshaler {
	state := &dualshock4.InputState{}

	if gp.GetButton(sdl.GamepadButtonSouth) {
		state.Buttons |= dualshock4.ButtonCross
	}
	if gp.GetButton(sdl.GamepadButtonEast) {
		state.Buttons |= dualshock4.ButtonCircle
	}
	if gp.GetButton(sdl.GamepadButtonWest) {
		state.Buttons |= dualshock4.ButtonSquare
	}
	if gp.GetButton(sdl.GamepadButtonNorth) {
		state.Buttons |= dualshock4.ButtonTriangle
	}
	if gp.GetButton(sdl.GamepadButtonLeftShoulder) {
		state.Buttons |= dualshock4.ButtonL1
	}
	if gp.GetButton(sdl.GamepadButtonRightShoulder) {
		state.Buttons |= dualshock4.ButtonR1
	}
	if gp.GetButton(sdl.GamepadButtonLeftStick) {
		state.Buttons |= dualshock4.ButtonL3
	}
	if gp.GetButton(sdl.GamepadButtonRightStick) {
		state.Buttons |= dualshock4.ButtonR3
	}
	if gp.GetButton(sdl.GamepadButtonStart) {
		state.Buttons |= dualshock4.ButtonOptions
	}
	if gp.GetButton(sdl.GamepadButtonBack) {
		state.Buttons |= dualshock4.ButtonShare
	}
	if gp.GetButton(sdl.GamepadButtonGuide) {
		state.Buttons |= dualshock4.ButtonPS
	}
	if gp.GetButton(sdl.GamepadButtonDpadUp) {
		state.DPad |= dualshock4.DPadUp
	}
	if gp.GetButton(sdl.GamepadButtonDpadDown) {
		state.DPad |= dualshock4.DPadDown
	}
	if gp.GetButton(sdl.GamepadButtonDpadLeft) {
		state.DPad |= dualshock4.DPadLeft
	}
	if gp.GetButton(sdl.GamepadButtonDpadRight) {
		state.DPad |= dualshock4.DPadRight
	}

	lt := gp.GetAxis(sdl.GamepadAxisLeftTrigger)
	rt := gp.GetAxis(sdl.GamepadAxisRightTrigger)

	state.L2 = uint8(max(0, min(math.MaxUint8, max(0, int32(lt))*math.MaxUint8/math.MaxInt16)))
	state.R2 = uint8(max(0, min(math.MaxUint8, max(0, int32(rt))*math.MaxUint8/math.MaxInt16)))
	if lt > 128 {
		state.Buttons |= dualshock4.ButtonL2
	}
	if rt > 128 {
		state.Buttons |= dualshock4.ButtonR2
	}

	state.LX = int8(max(math.MinInt8, min(math.MaxInt8, int32(gp.GetAxis(sdl.GamepadAxisLeftX))*(math.MaxInt8+1)/math.MaxInt16)))
	state.LY = int8(max(math.MinInt8, min(math.MaxInt8, int32(gp.GetAxis(sdl.GamepadAxisLeftY))*(math.MaxInt8+1)/math.MaxInt16)))
	state.RX = int8(max(math.MinInt8, min(math.MaxInt8, int32(gp.GetAxis(sdl.GamepadAxisRightX))*(math.MaxInt8+1)/math.MaxInt16)))
	state.RY = int8(max(math.MinInt8, min(math.MaxInt8, int32(gp.GetAxis(sdl.GamepadAxisRightY))*(math.MaxInt8+1)/math.MaxInt16)))

	return state
}
