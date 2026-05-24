package viiperdevice

import (
	"context"
	"encoding"
	"log/slog"
	"sync"

	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/VIIPER/apiclient"
	"github.com/Alia5/VIIPER/apitypes"
)

const stateBufferSize = 10

type DeviceType string

const (
	DeviceTypeUnknown    DeviceType = "unknown"
	DeviceTypeXbox360    DeviceType = "xbox360"
	DeviceTypeDualShock4 DeviceType = "dualshock4"
	DeviceTypeKeyboard   DeviceType = "keyboard"
)

type Device struct {
	controlStream *apiclient.DeviceStream
	deviceInfo    *apitypes.Device

	closeFunc func() error
	closeOnce sync.Once

	DeviceCtx       context.Context
	deviceCtxCancel context.CancelFunc

	FeedbackCh    <-chan encoding.BinaryUnmarshaler
	FeedbackErrCh <-chan error

	stateChan chan encoding.BinaryMarshaler
	done      chan struct{}
}

func New(
	controlStream *apiclient.DeviceStream,
	deviceInfo *apitypes.Device,
	closeFunc func() error,
) *Device {
	stateChan := make(chan encoding.BinaryMarshaler, stateBufferSize)
	deviceCtx, cancel := context.WithCancel(context.Background())

	decodeFeedback := readUnknownFeedback
	switch DeviceType(deviceInfo.Type) {
	case DeviceTypeKeyboard:
		decodeFeedback = readKeyboardFeedback
	case DeviceTypeDualShock4:
		decodeFeedback = readDualShock4Feedback
	case DeviceTypeXbox360:
		decodeFeedback = readXbox360Feedback
	}

	feedbackCh, errCh := controlStream.StartReading(deviceCtx, stateBufferSize, decodeFeedback)

	d := &Device{
		controlStream:   controlStream,
		deviceInfo:      deviceInfo,
		closeFunc:       closeFunc,
		stateChan:       stateChan,
		DeviceCtx:       deviceCtx,
		deviceCtxCancel: cancel,
		FeedbackCh:      feedbackCh,
		FeedbackErrCh:   errCh,
		done:            make(chan struct{}),
	}
	go d.handleState()
	return d
}

func (d *Device) Info() apitypes.Device {
	return *d.deviceInfo
}

func (d *Device) Update(gp *sdl.Gamepad) {
	if gp == nil {
		slog.Warn("Attempted to update VIIPER device with nil gamepad")
		return
	}

	var state encoding.BinaryMarshaler
	switch DeviceType(d.deviceInfo.Type) {
	case DeviceTypeXbox360:
		state = toXbox360State(gp)
	case DeviceTypeDualShock4:
		state = toDualShock4State(gp)
	// case DeviceTypeKeyboard:
	// 	state = toKeyboardState(gp)
	default:
		slog.Warn("Cant update unknown VIIPER device type", "device_type", d.deviceInfo.Type)
	}

	select {
	case <-d.done:
		return
	case d.stateChan <- state:
	default:
		slog.Warn("Dropping VIIPER device state update because buffer is full")
	}
}

func (d *Device) handleState() {
	for {
		select {
		case <-d.done:
			return
		case state := <-d.stateChan:
			err := d.controlStream.WriteBinary(state)
			if err != nil {
				slog.Error("Failed to send state to VIIPER device", "error", err)
				d.Close() //nolint
				return
			}
		}
	}
}

func (d *Device) IsClosed() bool {
	select {
	case <-d.done:
		return true
	default:
		return false
	}
}

func (d *Device) Close() error {
	var err error
	d.closeOnce.Do(func() {
		close(d.done)
		d.deviceCtxCancel()

		if d.controlStream != nil {
			err = d.controlStream.Close()
			d.controlStream = nil
		}

		if d.closeFunc != nil {
			go func(closeFunc func() error) {
				err := closeFunc()
				if err != nil {
					slog.Error("Failed to remove VIIPER device", "error", err)
				}
			}(d.closeFunc)
		}
	})

	return err
}
