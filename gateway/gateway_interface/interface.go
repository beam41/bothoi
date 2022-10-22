package gateway_interface

import (
	"bothoi/models/discord_models"
	"bothoi/models/types"
)

type ClientConnector interface {
	Connect()
}

type ClientVoiceConnector interface {
	JoinVoiceChannelMsg(guildId, channelId types.Snowflake, sessionIdChan chan<- string, voiceServerChan chan<- *discord_models.VoiceServer) error
	LeaveVoiceChannelMsg(guildId types.Snowflake) error
	CleanVoiceInstantiateChan(guildId types.Snowflake)
}

type ClientInterface interface {
	ClientConnector
	ClientVoiceConnector
}
