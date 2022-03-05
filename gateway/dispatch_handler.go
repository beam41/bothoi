package gateway

import (
	"bothoi/app_command"
	"bothoi/models"
	"bothoi/states"
	"log"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

func dispatchHandler(c *websocket.Conn, payload models.GatewayPayload) {
	switch payload.T {
	case "READY":
		mapstructure.Decode(payload.D, states.SessionState)
		states.SessionStateReady.Done()
	case "INTERACTION_CREATE":
		states.SessionStateReady.Wait()
		var data models.Interaction
		mapstructure.Decode(payload.D, &data)
		app_command.MapInteractionExecute(&data, c)
	case "GUILD_CREATE":
		states.SessionStateReady.Wait()
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
		states.SessionStateReady.Wait()
		var data *models.VoiceState = new(models.VoiceState)
		mapstructure.Decode(payload.D, data)
		if data.UserID != states.SessionState.User.ID {
			states.AddVoiceState(data)
		} else {
			states.SetSessionId(data.GuildID, data.SessionID)
			if states.SongQueue[data.GuildID].VoiceServer != nil {
				StartVoiceClient(data.GuildID)
			}
		}
	case "VOICE_SERVER_UPDATE":
		states.SessionStateReady.Wait()
		var data models.VoiceServer
		mapstructure.Decode(payload.D, &data)
		states.SetVoiceServer(data.GuildID, &data)
		if states.SongQueue[data.GuildID].SessionID != nil {
			StartVoiceClient(data.GuildID)
		}
	case "GUILD_UPDATE":
		// not important now
	case "GUILD_DELETE":
		// not important now
	}
}

func StartVoiceClient(guildId string) {
	log.Println("😘")
}
