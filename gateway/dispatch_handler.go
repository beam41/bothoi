package gateway

import (
	"bothoi/app_command"
	"bothoi/config"
	"bothoi/models"
	"bothoi/states"
	"bothoi/voice"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

func dispatchHandler(c *websocket.Conn, payload models.GatewayPayload) {
	switch payload.T {
	case "READY":
		var sessionState models.SessionState
		mapstructure.Decode(payload.D, &sessionState)
		states.AddSessionState(&sessionState)
	case "INTERACTION_CREATE":
		var data models.Interaction
		mapstructure.Decode(payload.D, &data)
		app_command.MapInteractionExecute(&data, c)
	case "GUILD_CREATE":
		var data models.Guild
		mapstructure.Decode(payload.D, &data)
		// guild voice state don't contain guild id
		states.AddGuild(&data)
		var voiceStates []models.VoiceState
		for _, voiceState := range data.VoiceStates {
			voiceState.GuildID = data.ID
			voiceStates = append(voiceStates, voiceState)
		}
		states.AddVoiceStateBulk(voiceStates)
	case "VOICE_STATE_UPDATE":
		var data *models.VoiceState = new(models.VoiceState)
		mapstructure.Decode(payload.D, data)
		if data.UserID != config.BOT_ID {
			states.AddVoiceState(data)
		} else {
			states.SetSessionId(data.GuildID, data.SessionID)
			if states.GetSongQueue(data.GuildID).VoiceServer != nil {
				StartVoiceClient(states.GetSongQueue(data.GuildID))
			}
		}
	case "VOICE_SERVER_UPDATE":
		var data models.VoiceServer
		mapstructure.Decode(payload.D, &data)
		states.SetVoiceServer(data.GuildID, &data)
		if states.GetSongQueue(data.GuildID).SessionID != nil {
			StartVoiceClient(states.GetSongQueue(data.GuildID))
		}
	case "GUILD_UPDATE":
		// not important now
	case "GUILD_DELETE":
		// not important now
	}
}

func StartVoiceClient(songQueue *models.SongQueue) {
	client := &voice.VoiceClient{
		SongQueue: songQueue,
	}
	go client.Connect()
}
