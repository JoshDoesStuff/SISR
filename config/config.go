package config

type UpdateNotify string

const (
	UpdateNotifyNone       UpdateNotify = "none"
	UpdateNotifyStable     UpdateNotify = "stable"
	UpdateNotifyPrerelease UpdateNotify = "prerelease"
)

type Log struct {
	Level string `help:"Log level: trace, debug, info, warn, error" aliases:"l" default:"info" env:"SISR_LOG_LEVEL"`
	File  string `help:"Log file path" env:"SISR_LOG_FILE"`
}

type Global struct {
	ConfigPath   string       `help:"Path to configuration file (json|yaml|toml)" name:"config" env:"SISR_CONFIG"`
	UpdateNotify UpdateNotify `help:"Update notification level: none, stable, prerelease" default:"stable" env:"SISR_UPDATE_NOTIFY"`
	Log          `embed:"" prefix:"log."`
	PlatformOpts `embed:""`
}

type API struct {
	ListenAddress   string `help:"API listen address" default:"localhost:0" env:"SISR_API_LISTEN_ADDRESS"`
	CORSOrigins     string `help:"CORS allowed origins" default:"https://steamloopback.host,http://steamloopback.host" env:"SISR_API_CORS_ORIGINS"`
	FrontendAddress string `help:"Frontend address (optional)" env:"SISR_API_FRONTEND_ADDRESS"`
}

type Steam struct {
	InstallDir         string `help:"Steam installation directory (optional, will attempt to auto-detect if not set)" env:"SISR_STEAM_INSTALL_DIR"`
	CEFRemoteDebugPort uint16 `help:"CEF remote debugging port (optional, will attempt to auto-detect if not set)" env:"SISR_STEAM_CEF_REMOTE_DEBUG_PORT"`
}
