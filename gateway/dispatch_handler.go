package gateway

import (
	"bothoi/commands"
	"bothoi/config"
	"bothoi/global"
	"bothoi/models/discord_models"
	"bothoi/repo"
	"bothoi/voice"
	"github.com/mitchellh/mapstructure"
	"log"
)

func mapInteractionExecute(data *discord_models.Interaction) {
	if interaction, ok := commands.ExecutorList[data.Data.Name]; ok {
		interaction(data)
	}
}

func dispatchHandler(payload discord_models.GatewayPayload) {
	switch payload.T {
	case "READY":
		var sessionState discord_models.ReadyEvent
		err := mapstructure.Decode(payload.D, &sessionState)
		if err != nil {
			log.Println(err)
			return
		}
		global.AddGatewaySession(&sessionState)
	case "INTERACTION_CREATE":
		var data discord_models.Interaction
		err := mapstructure.Decode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
		}
		mapInteractionExecute(&data)
	case "GUILD_CREATE":
		var data discord_models.GuildCreate
		err := mapstructure.Decode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
		}
		// guild voice state don't contain guild id
		repo.AddGuild(&data)
		var voiceStates []discord_models.VoiceState
		for _, voiceState := range data.VoiceStates {
			voiceState.GuildId = data.Id
			voiceStates = append(voiceStates, voiceState)
		}
		repo.AddVoiceStateBulk(voiceStates)
	case "VOICE_STATE_UPDATE":
		var data = new(discord_models.VoiceState)
		err := mapstructure.Decode(payload.D, data)
		if err != nil {
			log.Println(err)
			return
		}
		if data.UserId != config.BotId {
			repo.AddVoiceState(data)
		} else {
			voice.ReturnSessionId(data.GuildId, data.SessionId)
		}
	case "VOICE_SERVER_UPDATE":
		var data discord_models.VoiceServer
		err := mapstructure.Decode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
		}
		voice.ReturnVoiceServer(data.GuildId, &data)
	case "GUILD_UPDATE":
		// not important now
	case "GUILD_DELETE":
		// not important now
	}
}
