package input

import (
	"github.com/Alia5/SISR/sdl"
)

type Device struct {
	realGamepad         *sdl.Gamepad
	steamVirtualGamepad *sdl.Gamepad
	viiperDevice        *viiperDevice
}

func (d *Device) Gamepad(id sdl.GamepadID) (gp *sdl.Gamepad, isVirtual bool, err error) {
	if d.realGamepad != nil && d.realGamepad.ID() == id {
		return d.realGamepad, false, nil
	}
	if d.steamVirtualGamepad != nil && d.steamVirtualGamepad.ID() == id {
		return d.steamVirtualGamepad, true, nil
	}
	return gp, isVirtual, ErrNoDeviceForID
}

func (d *Device) UpdateViiperDevice(id sdl.GamepadID) error {
	gp, isVirtual, err := d.Gamepad(id)
	if err != nil {
		return err
	}
	if !isVirtual {
		return nil
	}

	d.viiperDevice.Update(gp)
	return nil
}
