package viiperdevice

import (
	"bufio"
	"encoding"
	"io"

	"github.com/Alia5/VIIPER/device/keyboard"
)

func readKeyboardFeedback(r *bufio.Reader) (encoding.BinaryUnmarshaler, error) {
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return nil, err
	}

	msg := new(keyboard.LEDState)
	if err := msg.UnmarshalBinary(b[:]); err != nil {
		return nil, err
	}

	return msg, nil
}
