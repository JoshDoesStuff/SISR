package steam

import "os"

func LaunchedViaSteam() (launchedViaSteam bool, launchedInGameMode bool) {

	steamAppIDEnv := os.Getenv("SteamAppId")
	launchedViaSteam = steamAppIDEnv != "" && steamAppIDEnv != "0"
	launchedInGameMode = os.Getenv("SteamInGameMode") != ""

	return launchedViaSteam, launchedInGameMode
}
