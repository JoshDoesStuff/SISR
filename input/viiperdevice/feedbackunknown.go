package viiperdevice

import (
	"bufio"
	"encoding"
	"io"
	"log/slog"
)

func readUnknownFeedback(_ *bufio.Reader) (encoding.BinaryUnmarshaler, error) {
	slog.Warn("Received feedback for device with unknown device type")
	return nil, io.EOF
}
