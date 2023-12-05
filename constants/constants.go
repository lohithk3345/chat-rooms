package constants

import "chat/internal/config"

var (
	HOST_URL      = config.EnvMap["HOST_URL"]
	HOST_PORT     = config.EnvMap["HOST_PORT"]
	WEBSOCKET_URL = config.EnvMap["WEBSOCKET_URL"]
)
