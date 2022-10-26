package gateway_interface

import (
	"bothoi/models/discord_models"
	"bothoi/models/types"
)

type ClientConnector interface {
	Connect()
}

type ClientVoiceConnector interface {
	JoinVoiceChannelMsg(guildID, channelID types.Snowflake, sessionIDChan chan<- string, voiceServerChan chan<- *discord_models.VoiceServer) error
	LeaveVoiceChannelMsg(guildID types.Snowflake) error
	CleanVoiceInstantiateChan(guildID types.Snowflake)
}

type ClientInterface interface {
	ClientConnector
	ClientVoiceConnector
}
