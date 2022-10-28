package gateway

import (
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"bothoi/repo"
	"github.com/mitchellh/mapstructure"
	"log"
)

func (client *Client) SetCommandExecutorList(executorList map[string]func(*Client, *discord_models.Interaction)) {
	client.commandExecutorList = executorList
}

func (client *Client) mapInteractionExecute(gatewayClient *Client, data *discord_models.Interaction) {
	if interaction, ok := client.commandExecutorList[data.Data.Name]; ok {
		interaction(gatewayClient, data)
	}
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
		client.mapInteractionExecute(client, &data)
	case "GUILD_CREATE":
		var data discord_models.GuildCreate
		err := mapstructure.WeakDecode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
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
			client.returnSessionID(data.GuildID, data.SessionID)
		}
	case "VOICE_SERVER_UPDATE":
		var data discord_models.VoiceServer
		err := mapstructure.WeakDecode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
		}
		client.returnVoiceServer(data.GuildID, &data)
	case "GUILD_UPDATE":
		// not important now
	case "GUILD_DELETE":
		// not important now
	}
}

func (client *Client) returnSessionID(guildID types.Snowflake, sessionID string) {
	client.voiceWaiter.RLock()
	defer client.voiceWaiter.RUnlock()
	if chanMap, ok := client.voiceWaiter.list[guildID]; ok {
		chanMap.sessionIDChan <- sessionID
	}
}

func (client *Client) returnVoiceServer(guildID types.Snowflake, voiceServer *discord_models.VoiceServer) {
	client.voiceWaiter.RLock()
	defer client.voiceWaiter.RUnlock()
	if chanMap, ok := client.voiceWaiter.list[guildID]; ok {
		chanMap.voiceServerChan <- voiceServer
	}
}
