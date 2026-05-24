package steam

import "errors"

var ErrOverlayLoadLaunchedViaSteam = errors.New("launched via Steam, overlay should already be loaded")
var ErrSteamNotRunning = errors.New("Steam is not running, loading overlay is useless") //nolint
var ErrMarkerNotFound = errors.New("SISR marker shortcut not found in Steam shortcuts") //nolint
var ErrShortcutsVDFNotFound = errors.New("shortcuts.vdf does not exist")
