package input

import "errors"

var ErrEmptyDevice = errors.New("device has no gamepads")
var ErrNoDeviceForID = errors.New("no device found for given id")

var ErrVirtualWithoutRealGamepad = errors.New("virtual gamepad without real gamepad")
var ErrVirtualAlreadyAssigned = errors.New("virtual gamepad already assigned to a real gamepad")
