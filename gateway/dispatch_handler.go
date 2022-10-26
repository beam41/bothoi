package gateway

import (
	"bothoi/bh_context"
	"bothoi/config"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"bothoi/repo"
	"github.com/mitchellh/mapstructure"
	"log"
)

func (client *client) dispatchHandler(payload discord_models.GatewayPayload) {
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
		bh_context.GetCommandManager().MapInteractionExecute(&data)
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

func (client *client) returnSessionID(guildID types.Snowflake, sessionID string) {
	client.voiceWaiter.RLock()
	defer client.voiceWaiter.RUnlock()
	if chanMap, ok := client.voiceWaiter.list[guildID]; ok {
		chanMap.sessionIDChan <- sessionID
	}
}

func (client *client) returnVoiceServer(guildID types.Snowflake, voiceServer *discord_models.VoiceServer) {
	client.voiceWaiter.RLock()
	defer client.voiceWaiter.RUnlock()
	if chanMap, ok := client.voiceWaiter.list[guildID]; ok {
		chanMap.voiceServerChan <- voiceServer
	}
}

func (client *client) JoinVoiceChannelMsg(guildID, channelID types.Snowflake, sessionIDChan chan<- string, voiceServerChan chan<- *discord_models.VoiceServer) error {
	createVoice := discord_models.NewVoiceStateUpdate(guildID, &channelID, false, true)
	err := client.gatewayConnWriteJSON(createVoice)
	if err != nil {
		return err
	}
	client.voiceWaiter.Lock()
	defer client.voiceWaiter.Unlock()
	client.voiceWaiter.list[guildID] = voiceInstantiateChan{sessionIDChan, voiceServerChan}
	return nil
}

func (client *client) LeaveVoiceChannelMsg(guildID types.Snowflake) error {
	client.CleanVoiceInstantiateChan(guildID)
	leaveVoice := discord_models.NewVoiceStateUpdate(guildID, nil, false, false)
	err := client.gatewayConnWriteJSON(leaveVoice)
	if err != nil {
		return err
	}
	return nil
}

func (client *client) CleanVoiceInstantiateChan(guildID types.Snowflake) {
	client.voiceWaiter.Lock()
	defer client.voiceWaiter.Unlock()
	delete(client.voiceWaiter.list, guildID)
}
