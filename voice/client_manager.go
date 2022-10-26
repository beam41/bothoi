package voice

import (
	"bothoi/bh_context"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"bothoi/repo"
	"bothoi/voice/voice_interface"
	"context"
	"errors"
	"log"
	"sync"
)

type clientManager struct {
	sync.RWMutex
	list map[types.Snowflake]*client
}

func NewClientManager() voice_interface.ClientManagerInterface {
	return &clientManager{
		list: make(map[types.Snowflake]*client),
	}
}

// createClient make new client from guild
func (clm *clientManager) createClient(guildID types.Snowflake) *client {
	if clm.list[guildID] != nil {
		return clm.list[guildID]
	}
	clm.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	clm.list[guildID] = &client{
		guildID:         guildID,
		udpReadyWait:    sync.NewCond(&sync.Mutex{}),
		ctx:             ctx,
		ctxCancel:       cancel,
		pauseWait:       sync.NewCond(&sync.Mutex{}),
		stopWaitForExit: make(chan struct{}),
		clm:             clm,
	}
	clm.Unlock()
	return clm.list[guildID]
}

// removeClient stop playing and remove the client from the list
func (clm *clientManager) removeClient(guildID types.Snowflake) error {
	clm.Lock()
	defer clm.Unlock()
	var client = clm.list[guildID]
	if client == nil {
		return errors.New("client not found")
	}
	client.Lock()
	defer client.Unlock()
	client.destroyed = true
	client.udpReadyWait.Broadcast()
	client.pauseWait.Broadcast()
	client.ctxCancel()
	client.closeConnection()
	delete(clm.list, guildID)
	return nil
}

// PauseClient pause/resume the music player return true if the player is paused
func (clm *clientManager) PauseClient(guildID types.Snowflake) (bool, error) {
	clm.Lock()
	defer clm.Unlock()
	var client = clm.list[guildID]
	if client == nil {
		return false, errors.New("client not found")
	}
	client.RLock()
	defer client.RUnlock()
	client.pauseWait.L.Lock()
	client.pausing = !client.pausing
	if !client.pausing {
		client.pauseWait.Broadcast()
	}
	client.pauseWait.L.Unlock()
	return client.pausing, nil
}

// SkipSong skip a song
func (clm *clientManager) SkipSong(guildID types.Snowflake) error {
	clm.Lock()
	defer clm.Unlock()
	var client = clm.list[guildID]
	if client == nil {
		return errors.New("client not found")
	}
	client.RLock()
	defer client.RUnlock()
	client.skip = true
	client.pauseWait.Broadcast()
	return nil
}

// StartClient start the client if not started already
func (clm *clientManager) StartClient(guildID, channelID types.Snowflake) error {
	clm.RLock()
	var client = clm.list[guildID]
	if client != nil {
		clm.RUnlock()
		client.RLock()
		defer client.RUnlock()
		if client.running {
			go client.play()
			return nil
		}
	}
	clm.RUnlock()

	log.Println(guildID, "Starting client")
	client = clm.createClient(guildID)
	sessionIDChan := make(chan string)
	voiceServerChan := make(chan *discord_models.VoiceServer)
	err := bh_context.GetGatewayClient().JoinVoiceChannelMsg(guildID, channelID, sessionIDChan, voiceServerChan)
	if err != nil {
		return err
	}

	// wait for session id and voice server
	go func() {
		sessionID, voiceServer := <-sessionIDChan, <-voiceServerChan
		bh_context.GetGatewayClient().CleanVoiceInstantiateChan(guildID)
		client.Lock()
		defer client.Unlock()
		client.sessionID = &sessionID
		client.voiceServer = voiceServer
		go client.connect()
		go client.play()
	}()

	return nil
}

// StopClient remove client from list and properly leave
func (clm *clientManager) StopClient(guildID types.Snowflake) error {
	_ = repo.DeleteSongsInGuild(guildID)
	err := clm.removeClient(guildID)
	if err != nil {
		return err
	}
	err = bh_context.GetGatewayClient().LeaveVoiceChannelMsg(guildID)
	if err != nil {
		return err
	}
	return nil
}
