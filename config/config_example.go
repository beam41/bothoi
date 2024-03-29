//go:build example
// +build example

package config

import (
	"bothoi/models/types"
	"time"
)

const (
	Development                                     = true
	BotToken                                        = ""
	BotID                           types.Snowflake = 0
	GatewayUrl                                      = "wss://gateway.discord.gg/?v=9&encoding=json"
	GatewayIntent                                   = 129
	InteractionResponseEndpoint                     = "https://discord.com/api/v8/interactions/<interaction_id>/<interaction_token>/callback"
	VoiceGatewayVersion                             = "4"
	PreferredMode                                   = "xsalsa20_poly1305"
	IdleTimeout                                     = time.Minute * 5
	DcaFramerate                                    = 48000
	DcaFrameduration                                = 20
	DcaBufferedFrame                                = int(time.Minute * 3 / (time.Millisecond * DcaFrameduration))
	InteractionResponseEditEndpoint                 = "https://discord.com/api/webhooks/<application_id>/<interaction_token>/messages/@original"
	CreateMessageEndpoint                           = "https://discord.com/api/channels/<channel_id>/messages"
)
