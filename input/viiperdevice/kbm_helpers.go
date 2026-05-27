package viiperdevice

import (
	"github.com/Alia5/VIIPER/device/keyboard"
	"github.com/Alia5/VIIPER/device/mouse"
)

func (d *Device) EnsureKeyboardState() *keyboard.InputState {
	state, ok := d.state.(*keyboard.InputState)
	if !ok || state == nil {
		state = &keyboard.InputState{}
		d.state = state
	}
	return state
}

func (d *Device) EnsureMouseState() *mouse.InputState {
	state, ok := d.state.(*mouse.InputState)
	if !ok || state == nil {
		state = &mouse.InputState{}
		d.state = state
	}
	return state
}

func (d *Device) ResetKBMState() {
	keyboardState := d.EnsureKeyboardState()
	keyboardState.Modifiers = 0
	clear(keyboardState.KeyBitmap[:])

	mouseState := d.EnsureMouseState()
	mouseState.Buttons = 0
	mouseState.DX = 0
	mouseState.DY = 0
	mouseState.Wheel = 0
	mouseState.Pan = 0
}
