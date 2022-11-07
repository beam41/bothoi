package gateway

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"bothoi/repo"
	"bothoi/util"
	"github.com/mitchellh/mapstructure"
	"log"
)

func (client *Client) SetInteractionExecutorList(executorList map[string]func(*discord_models.Interaction)) {
	client.interactionExecutorList = executorList
}

func (client *Client) interactionExecute(data *discord_models.Interaction) {
	if interaction, ok := client.interactionExecutorList[data.Data.Name]; ok {
		interaction(data)
	}
}

func (client *Client) RegisterNewSessionIDHandler(handler func(types.Snowflake, string)) {
	client.newSessionIDHandler = handler
}

func (client *Client) dispatchHandler(payload discord_models.GatewayPayload) {
	switch payload.T {
	case "READY":
		var sessionState discord_models.ReadyEvent
		err := mapstructure.WeakDecode(payload.D, &sessionState)
		if err != nil {
			log.Println(err)
			return
		}
		client.info.Lock()
		client.info.session = &sessionState
		client.info.Unlock()
	case "INTERACTION_CREATE":
		var data discord_models.Interaction
		err := mapstructure.WeakDecode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
		}
		client.interactionExecute(&data)
	case "GUILD_CREATE":
		var data discord_models.GuildCreate
		err := mapstructure.WeakDecode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
		}
		found := util.Find(data.VoiceStates, func(_ int, v discord_models.VoiceState) bool { return v.UserID == config.BotID })
		if found != nil {
			client.newSessionIDHandler(found.GuildID, found.SessionID)
		}
		repo.UpsertGuild(&data)
	case "VOICE_STATE_UPDATE":
		var data = new(discord_models.VoiceState)
		err := mapstructure.WeakDecode(payload.D, data)
		if err != nil {
			log.Println(err)
			return
		}
		repo.UpsertVoiceState(data)
		if data.UserID == config.BotID {
			if data.ChannelID == nil {
				return
			}
			client.voiceInstantiateList.RLock()
			defer client.voiceInstantiateList.RUnlock()
			if chanMap, ok := client.voiceInstantiateList.list[data.GuildID]; ok {
				chanMap.sessionIDChan <- data.SessionID
			} else {
				client.newSessionIDHandler(data.GuildID, data.SessionID)
			}
		}
	case "VOICE_SERVER_UPDATE":
		var data discord_models.VoiceServer
		err := mapstructure.WeakDecode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
		}
		client.voiceInstantiateList.RLock()
		defer client.voiceInstantiateList.RUnlock()
		if chanMap, ok := client.voiceInstantiateList.list[data.GuildID]; ok {
			chanMap.voiceServerChan <- &data
		}
	case "GUILD_UPDATE":
		// not important now
	case "GUILD_DELETE":
		// not important now
	}
}
