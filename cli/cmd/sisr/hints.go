package sisr

import (
	"log/slog"

	"github.com/Alia5/SISR/sdl"
)

var sdlHints = map[sdl.Hint]string{
	sdl.HintJoystickAllowBackgroundEvents: "1",
	sdl.HintHIDAPIIgnoreDevices:           "",
	// sdl.HintGamecontrollerIgnoreDevices:   "", // Set by env, this hint is blocked by Steam, but not from env...
	sdl.HintHIDAPIUdev:    "1",
	sdl.HintXInputEnabled: "1",
	// sdl.HintJoystickRawInput:           "1",
	// sdl.HintJoystickRawInputCorrelateXInput: "1",
	sdl.HintJoystickHIDAPI:                "1",
	sdl.HintJoystickHIDAPISwitch:          "1",
	sdl.HintJoystickHIDAPISwitch2:         "1",
	sdl.HintJoystickHIDAPIJoyCons:         "1",
	sdl.HintJoystickHIDAPINintendoClassic: "1",
	// sdl.HintJoystickHIDAPICombineJoyCons: "1",
	sdl.HintJoystickHIDAPIXbox:      "0", // TODO: rust code crashed when enabled, check if we can re-enable
	sdl.HintJoystickHIDAPIXbox360:   "0", // TODO: rust code crashed when enabled, check if we can re-enable
	sdl.HintJoystickHIDAPIXboxOne:   "0", // TODO: rust code crashed when enabled, check if we can re-enable
	sdl.HintJoystickEnhancedReports: "1",
	// sdl.HintJoystickDirectInput:       "1",
	// TODO: check if we can circumvent xbox-compatibility stuff
	// (eg. not being able to detect real/emulated controller for xbox controllers) with GameInput
	// sdl.HintJoystickGameInput:         "1",
	// sdl.HintHIDAPILibUSB:             "1",
}

func setSDLHints() {
	for hint, value := range sdlHints {
		err := sdl.SetHint(hint, value)
		if err != nil {
			slog.Error("Failed to set SDL hint", "hint", hint, "value", value, "error", err)
		}
	}
}
