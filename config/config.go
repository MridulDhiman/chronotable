package config

const (
	SNAPSHOT_EXT    = ".snap"
	MAIN_AOF_FILE   = "main.aof"
	CURR_AOF_FILE   = "curr.aof"
	CHRONO_MAIN_DIR = ".chrono"
	CONFIG_FILE     = "vulcan.toml"
)

var (
	ConfigKeyCurrVersion = "curr_version"
	ConfigKeyLatestVersion = "latest_version" 
)
