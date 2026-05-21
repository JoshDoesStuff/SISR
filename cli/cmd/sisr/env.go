package sisr

import "os"

func setEnv() {
	os.Setenv("SteamStreamingVideo", "0")                            //nolint:errcheck
	os.Setenv("SteamStreaming", "0")                                 //nolint:errcheck
	os.Setenv("SDL_GAMECONTROLLER_ALLOW_STEAM_VIRTUAL_GAMEPAD", "1") //nolint:errcheck
	os.Setenv("SDL_JOYSTICK_HIDAPI_STEAMXBOX", "1")                  //nolint:errcheck
	os.Setenv("SDL_GAMECONTROLLER_IGNORE_DEVICES", "")               //nolint:errcheck
	os.Setenv("SDL_GAMECONTROLLER_IGNORE_DEVICES_EXCEPT", "")        //nolint:errcheck
}
