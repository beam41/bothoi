package gateway

import (
	"bothoi/models/discord_models"
	"bothoi/models/types"
)

// PauseVoiceClient pause/resume the music player return true if the player is paused
func (client *Client) PauseVoiceClient(guildID types.Snowflake) (bool, error) {
	return client.voiceClientManager.PauseClient(guildID)
}

// SkipSong skip a song
func (client *Client) SkipSong(guildID types.Snowflake) error {
	return client.voiceClientManager.SkipSong(guildID)
}

// StartVoiceClient start the client if not started already
func (client *Client) StartVoiceClient(guildID, channelID types.Snowflake) error {
	connect, sessionIDChan, voiceServerChan := client.voiceClientManager.StartClient(guildID)
	if connect {
		createVoice := discord_models.NewVoiceStateUpdate(guildID, &channelID, false, true)
		err := client.gatewayConnWriteJSON(createVoice)
		if err != nil {
			return err
		}
		sessionIDChanRelay := make(chan string)
		voiceServerChanRelay := make(chan *discord_models.VoiceServer)
		client.voiceWaiter.Lock()
		client.voiceWaiter.list[guildID] = voiceInstantiateChan{sessionIDChanRelay, voiceServerChanRelay}
		client.voiceWaiter.Unlock()
		sessionIDChan <- <-sessionIDChanRelay
		voiceServerChan <- <-voiceServerChanRelay
		client.cleanVoiceInstantiateChan(guildID)
	}
	return nil
}

// StopVoiceClient remove client from list and properly leave
func (client *Client) StopVoiceClient(guildID types.Snowflake) error {
	err := client.voiceClientManager.StopClient(guildID)
	if err != nil {
		return err
	}
	client.cleanVoiceInstantiateChan(guildID)
	leaveVoice := discord_models.NewVoiceStateUpdate(guildID, nil, false, false)
	err = client.gatewayConnWriteJSON(leaveVoice)
	return err
}

func (client *Client) cleanVoiceInstantiateChan(guildID types.Snowflake) {
	client.voiceWaiter.Lock()
	defer client.voiceWaiter.Unlock()
	delete(client.voiceWaiter.list, guildID)
}
