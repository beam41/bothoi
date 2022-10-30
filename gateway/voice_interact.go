package gateway

import (
	"bothoi/models/discord_models"
	"bothoi/models/types"
)

// VoiceClientStart start the client if not started already
func (client *Client) VoiceChannelJoin(guildID, channelID types.Snowflake, sessionIDChan chan<- string, voiceServerChan chan<- *discord_models.VoiceServer) error {
	createVoice := discord_models.NewVoiceStateUpdate(guildID, &channelID, false, true)
	err := client.gatewayConnWriteJSON(createVoice)
	if err != nil {
		return err
	}
	client.voiceInstantiateList.Lock()
	defer client.voiceInstantiateList.Unlock()
	client.voiceInstantiateList.list[guildID] = voiceInstantiateChan{sessionIDChan, voiceServerChan}
	return nil
}

func (client *Client) VoiceChannelLeave(guildID types.Snowflake) error {
	client.CleanVoiceInstantiateChan(guildID)
	leaveVoice := discord_models.NewVoiceStateUpdate(guildID, nil, false, false)
	return client.gatewayConnWriteJSON(leaveVoice)
}

func (client *Client) CleanVoiceInstantiateChan(guildID types.Snowflake) {
	client.voiceInstantiateList.Lock()
	defer client.voiceInstantiateList.Unlock()
	delete(client.voiceInstantiateList.list, guildID)
}
