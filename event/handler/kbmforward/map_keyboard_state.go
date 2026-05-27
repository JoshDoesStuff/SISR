package kbmforward

import (
	"math"

	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/VIIPER/device/keyboard"
)

func mapKeyboardStateDown(state *keyboard.InputState, ev *sdl.KeyboardEvent) {
	if ev == nil || state == nil {
		return
	}

	if bit, ok := keyboardModifierBit(ev.Key); ok {
		state.Modifiers |= bit
		return
	}

	sc := ev.Scancode
	if sc > math.MaxUint8 {
		return
	}
	byteIdx := sc / 8
	bitIdx := uint(sc % 8)
	state.KeyBitmap[byteIdx] |= 1 << bitIdx
}

func mapKeyboardStateUp(state *keyboard.InputState, ev *sdl.KeyboardEvent) {
	if ev == nil || state == nil {
		return
	}

	if bit, ok := keyboardModifierBit(ev.Key); ok {
		state.Modifiers &^= bit
		return
	}

	sc := ev.Scancode
	if sc > math.MaxUint8 {
		return
	}
	byteIdx := sc / 8
	bitIdx := uint(sc % 8)
	state.KeyBitmap[byteIdx] &^= 1 << bitIdx
}

func keyboardModifierBit(key sdl.KeyCode) (uint8, bool) {
	switch key {
	case sdl.KeyCodeLCtrl:
		return keyboard.ModLeftCtrl, true
	case sdl.KeyCodeLShift:
		return keyboard.ModLeftShift, true
	case sdl.KeyCodeLAlt:
		return keyboard.ModLeftAlt, true
	case sdl.KeyCodeLGUI:
		return keyboard.ModLeftGUI, true
	case sdl.KeyCodeRCtrl:
		return keyboard.ModRightCtrl, true
	case sdl.KeyCodeRShift:
		return keyboard.ModRightShift, true
	case sdl.KeyCodeRAlt:
		return keyboard.ModRightAlt, true
	case sdl.KeyCodeRGUI:
		return keyboard.ModRightGUI, true
	default:
		return 0, false
	}
}
