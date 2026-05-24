package viiperdevice

import (
	"bufio"
	"encoding"
	"io"

	"github.com/Alia5/VIIPER/device/xbox360"
)

func readXbox360Feedback(r *bufio.Reader) (encoding.BinaryUnmarshaler, error) {
	var b [2]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return nil, err
	}

	msg := new(xbox360.XRumbleState)
	if err := msg.UnmarshalBinary(b[:]); err != nil {
		return nil, err
	}

	return msg, nil
}
