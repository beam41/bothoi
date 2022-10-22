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
		err := mapstructure.Decode(payload.D, &sessionState)
		if err != nil {
			log.Println(err)
			return
		}
		client.info.Lock()
		client.info.session = &sessionState
		client.info.Unlock()
	case "INTERACTION_CREATE":
		var data discord_models.Interaction
		err := mapstructure.Decode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
		}
		bh_context.GetCommandManager().MapInteractionExecute(&data)
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
			client.returnSessionId(data.GuildId, data.SessionId)
		}
	case "VOICE_SERVER_UPDATE":
		var data discord_models.VoiceServer
		err := mapstructure.Decode(payload.D, &data)
		if err != nil {
			log.Println(err)
			return
		}
		client.returnVoiceServer(data.GuildId, &data)
	case "GUILD_UPDATE":
		// not important now
	case "GUILD_DELETE":
		// not important now
	}
}

func (client *client) returnSessionId(guildId types.Snowflake, sessionId string) {
	client.voiceWaiter.RLock()
	defer client.voiceWaiter.RUnlock()
	if chanMap, ok := client.voiceWaiter.list[guildId]; ok {
		chanMap.sessionIdChan <- sessionId
	}
}

func (client *client) returnVoiceServer(guildId types.Snowflake, voiceServer *discord_models.VoiceServer) {
	client.voiceWaiter.RLock()
	defer client.voiceWaiter.RUnlock()
	if chanMap, ok := client.voiceWaiter.list[guildId]; ok {
		chanMap.voiceServerChan <- voiceServer
	}
}

func (client *client) JoinVoiceChannelMsg(guildId, channelId types.Snowflake, sessionIdChan chan<- string, voiceServerChan chan<- *discord_models.VoiceServer) error {
	createVoice := discord_models.NewVoiceStateUpdate(guildId, &channelId, false, true)
	err := client.gatewayConnWriteJSON(createVoice)
	if err != nil {
		return err
	}
	client.voiceWaiter.Lock()
	defer client.voiceWaiter.Unlock()
	client.voiceWaiter.list[guildId] = voiceInstantiateChan{sessionIdChan, voiceServerChan}
	return nil
}

func (client *client) LeaveVoiceChannelMsg(guildId types.Snowflake) error {
	client.CleanVoiceInstantiateChan(guildId)
	leaveVoice := discord_models.NewVoiceStateUpdate(guildId, nil, false, false)
	err := client.gatewayConnWriteJSON(leaveVoice)
	if err != nil {
		return err
	}
	return nil
}

func (client *client) CleanVoiceInstantiateChan(guildId types.Snowflake) {
	client.voiceWaiter.Lock()
	defer client.voiceWaiter.Unlock()
	delete(client.voiceWaiter.list, guildId)
}
