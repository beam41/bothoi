//go:build example
// +build example

package config

const (
	NO_COMMAND_REGISTER                = true
	DEVELOPMENT                        = true
	BOT_TOKEN                          = ""
	BOT_ID                             = ""
	GATEWAY_URL                        = "wss://gateway.discord.gg/?v=9&encoding=json"
	GATEWAY_INTENT                     = 129
	APP_COMMAND_ENDPOINT               = ""
	INTERACTION_RESPONSE_ENDPOINT      = "https://discord.com/api/v8/interactions/<interaction_id>/<interaction_token>/callback"
	APP_COMMAND_GUILD_ENDPOINT         = ""
	VOICE_GATEWAY_VERSION              = "4"
	PREFERRED_MODE                     = "xsalsa20_poly1305"
	UDP_TIMEOUT                        = time.Millisecond * 500
	IDLE_TIMEOUT                       = time.Minute * 5
	DCA_FRAMERATE                      = 48000
	DCA_FRAMEDURATION                  = 20
	INTERACTION_RESPONSE_EDIT_ENDPOINT = "https://discord.com/api/webhooks/<application_id>/<interaction_token>/messages/@original"
)
