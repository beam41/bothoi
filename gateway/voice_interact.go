package gateway

import (
	"bothoi/models/discord_models"
	"bothoi/models/types"
)

// VoiceClientStart start the client if not started already
func (client *Client) VoiceClientStart(guildID, channelID types.Snowflake) error {
	connect, sessionIDChan, voiceServerChan := client.voiceClientManager.ClientStart(guildID)
	if connect {
		createVoice := discord_models.NewVoiceStateUpdate(guildID, &channelID, false, true)
		err := client.gatewayConnWriteJSON(createVoice)
		if err != nil {
			return err
		}
		sessionIDChanRelay := make(chan string)
		voiceServerChanRelay := make(chan *discord_models.VoiceServer)
		client.voiceInstantiateList.Lock()
		client.voiceInstantiateList.list[guildID] = voiceInstantiateChan{sessionIDChanRelay, voiceServerChanRelay}
		client.voiceInstantiateList.Unlock()
		sessionIDChan <- <-sessionIDChanRelay
		voiceServerChan <- <-voiceServerChanRelay
		client.cleanVoiceInstantiateChan(guildID)
	}
	return nil
}

// VoiceClientStop remove client from list and properly leave
func (client *Client) VoiceClientStop(guildID types.Snowflake) (success bool, err error) {
	success = client.voiceClientManager.ClientStop(guildID)
	if !success {
		return
	}
	client.cleanVoiceInstantiateChan(guildID)
	leaveVoice := discord_models.NewVoiceStateUpdate(guildID, nil, false, false)
	err = client.gatewayConnWriteJSON(leaveVoice)
	if err != nil {
		success = false
	}
	return
}

func (client *Client) cleanVoiceInstantiateChan(guildID types.Snowflake) {
	client.voiceInstantiateList.Lock()
	defer client.voiceInstantiateList.Unlock()
	delete(client.voiceInstantiateList.list, guildID)
}

// VoiceClientPauseSong pause/resume the music player return true if the player is paused
func (client *Client) VoiceClientPauseSong(guildID types.Snowflake) (found bool, pausing bool) {
	return client.voiceClientManager.ClientPauseSong(guildID)
}

// VoiceClientSkipSong skip a song
func (client *Client) VoiceClientSkipSong(guildID types.Snowflake) bool {
	return client.voiceClientManager.ClientSkipSong(guildID)
}
