package gateway

import (
	"bothoi/commands"
	"bothoi/config"
	"bothoi/models"
	"bothoi/repo"
	"bothoi/voice"
	"github.com/mitchellh/mapstructure"
	"log"
)

func mapInteractionExecute(data *models.Interaction) {
	if interaction, ok := commands.ExecutorList[data.Data.Name]; ok {
		interaction(data)
	}
}

func dispatchHandler(payload models.GatewayPayload) {
	switch payload.T {
	case "READY":
		var sessionState models.SessionState
		err := mapstructure.Decode(payload.D, &sessionState)
		if err != nil {
			log.Println(err)
			return
		}
		repo.AddSessionState(&sessionState)
	case "INTERACTION_CREATE":
		var data models.Interaction
		err := mapstructure.Decode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
		}
		mapInteractionExecute(&data)
	case "GUILD_CREATE":
		var data models.Guild
		err := mapstructure.Decode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
		}
		// guild voice state don't contain guild id
		repo.AddGuild(&data)
		var voiceStates []models.VoiceState
		for _, voiceState := range data.VoiceStates {
			voiceState.GuildID = data.ID
			voiceStates = append(voiceStates, voiceState)
		}
		repo.AddVoiceStateBulk(voiceStates)
	case "VOICE_STATE_UPDATE":
		var data = new(models.VoiceState)
		err := mapstructure.Decode(payload.D, data)
		if err != nil {
			log.Println(err)
			return
		}
		if data.UserID != config.BotId {
			repo.AddVoiceState(data)
		} else {
			voice.ReturnSessionId(data.GuildID, data.SessionID)
		}
	case "VOICE_SERVER_UPDATE":
		var data models.VoiceServer
		err := mapstructure.Decode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
		}
		voice.ReturnVoiceServer(data.GuildID, &data)
	case "GUILD_UPDATE":
		// not important now
	case "GUILD_DELETE":
		// not important now
	}
}
