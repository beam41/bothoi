//go:build example
// +build example

package config

import "time"

const (
	Development                     = true
	BotToken                        = ""
	BotId                           = ""
	GatewayUrl                      = "wss://gateway.discord.gg/?v=9&encoding=json"
	GatewayIntent                   = 129
	InteractionResponseEndpoint     = "https://discord.com/api/v8/interactions/<interaction_id>/<interaction_token>/callback"
	VoiceGatewayVersion             = "4"
	PreferredMode                   = "xsalsa20_poly1305"
	IdleTimeout                     = time.Minute * 5
	DcaFramerate                    = 48000
	DcaFrameduration                = 20
	InteractionResponseEditEndpoint = "https://discord.com/api/webhooks/<application_id>/<interaction_token>/messages/@original"
)
