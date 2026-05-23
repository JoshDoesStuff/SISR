package input

import (
	"log/slog"
	"sync"

	"github.com/Alia5/SISR/sdl"
)

type Device struct {
	RealGamepad         *sdl.Gamepad
	SteamVirtualGamepad *sdl.Gamepad
	ViiperDevice        *ViiperDevice

	mtx sync.Mutex
}

func (d *Device) Lock() {
	d.mtx.Lock()
}

func (d *Device) Unlock() {
	d.mtx.Unlock()
}

func (d *Device) Gamepad(id sdl.GamepadID) (gp *sdl.Gamepad, isVirtual bool, err error) {
	if d.RealGamepad != nil && d.RealGamepad.ID() == id {
		return d.RealGamepad, false, nil
	}
	if d.SteamVirtualGamepad != nil && d.SteamVirtualGamepad.ID() == id {
		return d.SteamVirtualGamepad, true, nil
	}
	return gp, isVirtual, ErrNoDeviceForID
}

func (d *Device) Close() {
	if d.RealGamepad != nil {
		d.RealGamepad.Close()
		d.RealGamepad = nil
	}
	if d.SteamVirtualGamepad != nil {
		d.SteamVirtualGamepad.Close()
		d.SteamVirtualGamepad = nil
	}
	if d.ViiperDevice != nil {
		err := d.ViiperDevice.Close()
		if err != nil {
			slog.Error("Failed to close VIIPER device", "error", err)
		}
		d.ViiperDevice = nil
	}
}

func (d *Device) SetViiperDevice(vd *ViiperDevice) {
	if d.ViiperDevice != nil {
		slog.Warn("Device already has a VIIPER device assigned, Overwriting")
		err := d.ViiperDevice.Close()
		if err != nil {
			slog.Error("Failed to close existing VIIPER device when setting new one", "error", err)
		}
	}
	d.ViiperDevice = vd
}
