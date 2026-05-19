package main

import (
	"log/slog"
	"os"
	"runtime"
	"slices"
	"time"

	"github.com/Alia5/SISR/logging"
	"github.com/Alia5/SISR/sdl"
	"github.com/Alia5/SISR/sdl/extras"
	"github.com/Alia5/SISR/webview"
	"github.com/Alia5/SISR/windows"
	"github.com/Alia5/SISR/windows/hooks"
)

func main() {

	if _, _, err := logging.SetupLogger("debug", ""); err != nil {
		slog.Error("Failed to setup logger", "error", err)
	}

	// The SDL_GAMECONTROLLER_IGNORE_DEVICES hint doesn't work when Steam is
	// injected, but the env-var form does.
	os.Setenv("SteamStreamingVideo", "0")                            //nolint:errcheck
	os.Setenv("SteamStreaming", "0")                                 //nolint:errcheck
	os.Setenv("SDL_GAMECONTROLLER_ALLOW_STEAM_VIRTUAL_GAMEPAD", "1") //nolint:errcheck
	os.Setenv("SDL_JOYSTICK_HIDAPI_STEAMXBOX", "1")                  //nolint:errcheck
	os.Setenv("SDL_GAMECONTROLLER_IGNORE_DEVICES", "")               //nolint:errcheck
	os.Setenv("SDL_GAMECONTROLLER_IGNORE_DEVICES_EXCEPT", "")        //nolint:errcheck

	err := sdl.Init(sdl.InitFlagVideo | sdl.InitFlagGamepad)
	if err != nil {
		slog.Error("Failed to init SDL", "error", err)
		os.Exit(1)
	}
	defer sdl.Quit()

	sdl.SetGamepadEventsEnabled(true)

	gamepads := map[sdl.GamepadID]*sdl.Gamepad{}
	defer func() {
		for id, gp := range gamepads {
			if gp != nil {
				gp.Close()
			}
			delete(gamepads, id)
		}
	}()

	openGamepad := func(id sdl.GamepadID) {
		if gp, ok := gamepads[id]; ok && gp != nil {
			return
		}
		gp, openErr := sdl.OpenGamepad(id)
		if openErr != nil {
			slog.Warn("Failed to open gamepad", "id", id, "error", openErr)
			return
		}
		gamepads[id] = gp
		slog.Info("Gamepad opened", "id", id, "name", gp.Name(), "type", gp.Type(), "steamhandle", gp.GetSteamHandle())
	}

	closeGamepad := func(id sdl.GamepadID) {
		gp, ok := gamepads[id]
		if !ok {
			return
		}
		if gp != nil {
			gp.Close()
		}
		delete(gamepads, id)
		slog.Info("Gamepad closed", "id", id)
	}

	if ids, listErr := sdl.GetGamepads(); listErr != nil {
		slog.Warn("Failed to list initial gamepads", "error", listErr)
	} else {
		for _, id := range ids {
			openGamepad(id)
		}
	}

	window, renderer, err := sdl.CreateWindowAndRenderer(
		"SISR",
		1280,
		720,
		sdl.WindowFlagVulkan|
			sdl.WindowFlagTransparent|
			sdl.WindowFlagBorderless|
			sdl.WindowFlagAlwaysOnTop,
	)
	if err != nil {
		slog.Error("Failed to create window", "error", err)
		os.Exit(1)
	}
	defer window.Destroy()
	defer renderer.Destroy()

	if runtime.GOOS == "windows" {
		err = windows.SetDWMPassiveUpdateMode(window)
		if err != nil {
			slog.Error("Failed to enable DWM passive update mode", "error", err)
		}
	}

	err = extras.SetCursorHitTest(window, false)
	if err != nil {
		slog.Error("Failed to set cursor hit test", "error", err)
	}

	err = renderer.SetRenderDrawColor(0, 0, 80, 128)
	if err != nil {
		slog.Error("Failed to set render draw color", "error", err)
		os.Exit(1)
	}

	wv, err := webview.New(window, 1280, 720, true)
	if err != nil {
		slog.Error("Failed to create webview", "error", err)
		os.Exit(1)
	}
	defer wv.Destroy()
	wv.SetHTML(`<!DOCTYPE html>
<html style="background: transparent;">
<body style="margin: 0; background: transparent; display: grid; place-items: center; height: 100svh;">
	<h1 style="color: red;">
		hello SISR ✂️
	</h1>
</body>
</html>`)

	if runtime.GOOS == "windows" {
		hookedFns := hooks.DetectHooks("hid.dll")
		if len(hookedFns) > 0 {
			slog.Info("Detected HID hooks")
			toUnhookFNs := []string{
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

			for _, toUnhook := range toUnhookFNs {
				if slices.Contains(hookedFns, toUnhook) {
					unhooked := hooks.Unhook("hid.dll", toUnhook)
					if unhooked {
						slog.Info("Successfully unhooked HID export", "export", toUnhook)
					} else {
						slog.Warn("Failed to unhook HID export", "export", toUnhook)
					}
				} else {
					slog.Info("HID export not hooked, skipping", "export", toUnhook)
				}
			}
		}
	}

	for {
		ev, _ := sdl.WaitEventTimeout(time.Millisecond * 16)
		if ev != nil {
			if runtime.GOOS == "linux" {
				err := extras.HandleCursorHitTestWindowEvent(window, ev)
				if err != nil {
					slog.Error("Failed to handle cursor hit test window event", "error", err)
				}
			}
			switch ev := ev.(type) {
			case *sdl.QuitEvent:
				return
			case *sdl.KeyboardEvent:
				if ev.Key == sdl.KeyCodeEscape && ev.Down {
					wv.SetVisible(!wv.Visible())
					if wv.Visible() {
						slog.Info("WebView shown")
					} else {
						slog.Info("WebView hidden")
					}
				}
			case *sdl.GamepadDeviceEvent:
				switch ev.Type {
				case sdl.EventTypeGamepadAdded:
					id := sdl.GamepadID(ev.Which)
					slog.Info("Gamepad connected", "id", id, "name", sdl.GetGamepadNameForID(id))
					openGamepad(id)
				case sdl.EventTypeGamepadRemoved:
					id := sdl.GamepadID(ev.Which)
					slog.Info("Gamepad disconnected", "id", id)
					closeGamepad(id)
				}
			case *sdl.WindowEvent:
				if ev.Type == sdl.EventTypeWindowResized {
					wv.Resize(int(ev.Data1), int(ev.Data2))
				}
			}
		}

		wv.Tick()
		err = renderer.RenderClear()
		if err != nil {
			slog.Error("Failed to clear renderer", "error", err)
			os.Exit(1)
		}
		err = renderer.RenderPresent()
		if err != nil {
			slog.Error("Failed to present renderer", "error", err)
			os.Exit(1)
		}
	}

}
