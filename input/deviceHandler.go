package input

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"slices"

	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/webview"
	"github.com/Alia5/SISR/windows/hooks"
)

type DeviceHandler interface {
	OpenGamePad(id sdl.GamepadID) (*Device, error)
	CloseGamePad(id sdl.GamepadID) error
	DeviceForID(id sdl.GamepadID) (*Device, bool)
}

var toUnhookFNs = []string{
	"HidD_FreePreparsedData",
	"HidD_GetAttributes",
	"HidD_GetPreparsedData",
	"HidD_GetProductString",
	"HidP_GetButtonCaps",
	"HidP_GetCaps",
	"HidP_GetData",
	"HidP_GetUsageValue",
	"HidP_GetUsages",
	"HidP_GetValueCaps",
	"HidP_MaxDataListLength",
}

type deviceHandler struct {
	window  *sdl.Window
	webview webview.WebView
	devices map[sdl.GamepadID]*Device
	viiper  ViiperBridge
}

func NewDeviceHandler(window *sdl.Window, webview webview.WebView) (DeviceHandler, func(), error) {

	if runtime.GOOS == "windows" {
		hookedFns := hooks.DetectHooks("hid.dll")
		if len(hookedFns) > 0 {
			slog.Info("Detected HID hooks")

			for _, toUnhook := range toUnhookFNs {
				if slices.Contains(hookedFns, toUnhook) {
					unhooked := hooks.Unhook("hid.dll", toUnhook)
					if unhooked {
						slog.Debug("Successfully unhooked HID export", "export", toUnhook)
					} else {
						slog.Warn("Failed to unhook HID export", "export", toUnhook)
					}
				} else {
					slog.Debug("HID export not hooked, skipping", "export", toUnhook)
				}
			}
		}
	}

	err := sdl.InitSubSystem(sdl.InitFlagGamepad | sdl.InitFlagSensor | sdl.InitFlagHaptic)
	if err != nil {
		return nil, nil, err
	}

	dh := &deviceHandler{
		window:  window,
		webview: webview,
		devices: make(map[sdl.GamepadID]*Device),
		viiper:  NewViiperBridge(), // TODO:
	}
	return dh, dh.quit, nil
}

func (dh *deviceHandler) OpenGamePad(id sdl.GamepadID) (*Device, error) {
	if dev, exists := dh.devices[id]; exists {
		return dev, nil
	}
	var dev *Device
	gp, err := sdl.OpenGamepad(id)
	if err != nil {
		return nil, err
	}
	steamHandle := gp.GetSteamHandle()
	serial := gp.Serial()
	path := gp.Path()
	gpType := gp.Type()
	realType := gp.RealType()

	slog.Debug("Opened Gamepad",
		"id", id,
		"name", gp.Name(),
		"steamHandle", steamHandle,
		"serial", serial,
		"path", path,
		"type", gpType.Name(),
		"realType", realType.Name(),
	)

	if steamHandle != 0 {
		if id == 0 {
			// Xbox controller, ignore for now
			gp.Close()
			return nil, fmt.Errorf(
				"%w: first encountered device id=%d (SteamHandle=%d) has no corresponding real gamepad",
				ErrVirtualWithoutRealGamepad, id, steamHandle,
			)
		}
		// We assume devices are opened in order
		// Real Gamepad first and Steam virtual gamepad immediately after with Id+1
		// I havent validated, but also haven't encountered anything else so far.
		var ok bool
		dev, ok = dh.devices[id-1]
		if !ok {
			// Xbox controller, ignore for now
			gp.Close()
			return nil, fmt.Errorf(
				"%w: device id=%d (SteamHandle=%d) has no corresponding real gamepad",
				ErrVirtualWithoutRealGamepad, id, steamHandle,
			)
		}
		if dev.steamVirtualGamepad != nil {
			gp.Close()
			return nil, ErrVirtualAlreadyAssigned
		}
		dev.steamVirtualGamepad = gp
		slog.Info("Found Steam Virtual Gamepad",
			"id", id,
			"name", gp.Name(),
			"steamhandle", steamHandle,
			"paired with", dev.realGamepad.ID(),
			"paired with name", dev.realGamepad.Name(),
		)
		// TODO:
		err := dh.viiper.AttachViiperDevice(context.TODO(), dev)
		if err != nil {
			slog.Error("Failed to attach VIIPER device", "error", err)
		}
	} else {
		dev = &Device{realGamepad: gp}
		slog.Debug("Opened real gamepad", "id", id, "name", gp.Name())
	}
	dh.devices[id] = dev

	return dev, nil
}

func (dh *deviceHandler) CloseGamePad(id sdl.GamepadID) error {
	dev, ok := dh.devices[id]
	if !ok {
		return nil
	}
	gp := dev.realGamepad
	if gp == nil {
		gp = dev.steamVirtualGamepad
		dev.steamVirtualGamepad = nil
	} else {
		dev.realGamepad = nil
	}
	slog.Debug("Closing gamepad", "id", id, "name", gp.Name(), "steamHandle", gp.GetSteamHandle())
	gp.Close()

	delete(dh.devices, id)
	return nil
}

func (dh *deviceHandler) DeviceForID(id sdl.GamepadID) (*Device, bool) {
	dev, ok := dh.devices[id]
	return dev, ok
}

func (dh *deviceHandler) UpdateDeviceState(id sdl.GamepadID) error {

	return nil
}

func (dh *deviceHandler) quit() {
	for id, dev := range dh.devices {
		if dev != nil {
			gp := dev.realGamepad
			if gp == nil {
				gp = dev.steamVirtualGamepad
				dev.steamVirtualGamepad = nil
			} else {
				dev.realGamepad = nil
			}
			slog.Debug("Closing gamepad", "id", id, "name", gp.Name(), "steamHandle", gp.GetSteamHandle())
			gp.Close()

			// TODO: remove VIIPER device

		}
		delete(dh.devices, id)
	}
}
