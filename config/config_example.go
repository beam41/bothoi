//go:build example
// +build example

package config

const (
	DEVELOPMENT                   = true
	BOT_TOKEN                     = ""
	GATEWAY_URL                   = "wss://gateway.discord.gg/?v=9&encoding=json"
	GATEWAY_INTENT                = 129
	APP_COMMAND_ENDPOINT          = ""
	INTERACTION_RESPONSE_ENDPOINT = "https://discord.com/api/v8/interactions/<interaction_id>/<interaction_token>/callback"
	APP_COMMAND_GUILD_ENDPOINT    = ""
)
