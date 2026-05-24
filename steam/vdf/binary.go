package vdf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

const (
	typeMapStart byte = 0x00
	typeString   byte = 0x01
	typeNumber   byte = 0x02
	typeMapEnd   byte = 0x08
)

var (
	ErrUnsupportedType = errors.New("unsupported type")
)

func Read(r io.Reader) (map[string]any, error) {
	return readMap(r)
}

func readMap(r io.Reader) (map[string]any, error) {
	out := make(map[string]any)

	for {
		typ, err := readByte(r)
		if err != nil {
			return nil, err
		}

		if typ == typeMapEnd {
			return out, nil
		}

		key, err := readString(r)
		if err != nil {
			return nil, err
		}

		value, err := readValue(r, typ)
		if err != nil {
			return nil, err
		}

		out[key] = value
	}
}

func readValue(r io.Reader, typ byte) (any, error) {
	switch typ {
	case typeMapStart:
		return readMap(r)
	case typeString:
		return readString(r)
	case typeNumber:
		return readUint32(r)
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedType, typ)
	}
}

func readByte(r io.Reader) (byte, error) {
	var b [1]byte
	_, err := io.ReadFull(r, b[:])
	if err != nil {
		return 0, err
	}

	return b[0], nil
}

func readString(r io.Reader) (string, error) {
	buf := make([]byte, 0, 32)

	for {
		b, err := readByte(r)
		if err != nil {
			return "", err
		}

		if b == 0 {
			return string(buf), nil
		}

		buf = append(buf, b)
	}
}

func readUint32(r io.Reader) (uint32, error) {
	var raw [4]byte
	_, err := io.ReadFull(r, raw[:])
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(raw[:]), nil
}
