package viiperdevice

import (
	"context"
	"encoding"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/VIIPER/apiclient"
	"github.com/Alia5/VIIPER/apitypes"
)

const stateBufferSize = 32

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

	state     encoding.BinaryMarshaler
	stateChan chan []byte
	done      chan struct{}
}

func New(
	controlStream *apiclient.DeviceStream,
	deviceInfo *apitypes.Device,
	closeFunc func() error,
) *Device {
	stateChan := make(chan []byte, stateBufferSize)
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

	switch DeviceType(d.deviceInfo.Type) {
	case DeviceTypeXbox360:
		toXbox360State(gp, &d.state)
	case DeviceTypeDualShock4:
		toDualShock4State(gp, &d.state)
	// case DeviceTypeKeyboard:
	// 	state = toKeyboardState(gp)
	default:
		slog.Warn("Cant update unknown VIIPER device type", "device_type", d.deviceInfo.Type)
		return
	}

	if d.state == nil {
		slog.Warn("No VIIPER state available to marshal", "device_type", d.deviceInfo.Type)
		return
	}

	data, err := d.state.MarshalBinary()
	if err != nil {
		slog.Error("Failed to marshal VIIPER device state", "error", err)
		return
	}

	timer := time.NewTimer(1 * time.Second)
	defer timer.Stop()

	select {
	case <-d.done:
		slog.Warn("Attempted to update VIIPER device after it was closed")
		return
	case d.stateChan <- data:
		// sent successfully
	case <-timer.C:
		slog.Warn("Timed out sending state update to VIIPER device;")
	}
}

func (d *Device) handleState() {
	stream := d.controlStream

	for {
		select {
		case <-d.done:
			return
		case data := <-d.stateChan:
			n, err := stream.Write(data)
			if err == nil && n != len(data) {
				err = fmt.Errorf("short write: wrote %d of %d bytes", n, len(data))
			}
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
