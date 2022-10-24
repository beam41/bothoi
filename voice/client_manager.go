package voice

import (
	"bothoi/bh_context"
	"bothoi/models/discord_models"
	"bothoi/models/types"
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
func (clm *clientManager) createClient(guildId types.Snowflake) *client {
	if clm.list[guildId] != nil {
		return clm.list[guildId]
	}
	clm.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	clm.list[guildId] = &client{
		guildId:         guildId,
		udpReadyWait:    sync.NewCond(&sync.Mutex{}),
		ctx:             ctx,
		ctxCancel:       cancel,
		pauseWait:       sync.NewCond(&sync.Mutex{}),
		stopWaitForExit: make(chan struct{}),
		clm:             clm,
	}
	clm.Unlock()
	return clm.list[guildId]
}

// removeClient stop playing and remove the client from the list
func (clm *clientManager) removeClient(guildId types.Snowflake) error {
	clm.Lock()
	defer clm.Unlock()
	var client = clm.list[guildId]
	if client == nil {
		return errors.New("client not found")
	}
	client.Lock()
	defer client.Unlock()
	client.destroyed = true
	client.udpReadyWait.Broadcast()
	client.pauseWait.Broadcast()
	client.ctxCancel()
	delete(clm.list, guildId)
	return nil
}

// PauseClient pause/resume the music player return true if the player is paused
func (clm *clientManager) PauseClient(guildId types.Snowflake) (bool, error) {
	clm.Lock()
	defer clm.Unlock()
	var client = clm.list[guildId]
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
func (clm *clientManager) SkipSong(guildId types.Snowflake) error {
	clm.Lock()
	defer clm.Unlock()
	var client = clm.list[guildId]
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
func (clm *clientManager) StartClient(guildId, channelId types.Snowflake) error {
	clm.RLock()
	var client = clm.list[guildId]
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

	log.Println("Starting client")
	client = clm.createClient(guildId)
	sessionIdChan := make(chan string)
	voiceServerChan := make(chan *discord_models.VoiceServer)
	err := bh_context.GetGatewayClient().JoinVoiceChannelMsg(guildId, channelId, sessionIdChan, voiceServerChan)
	if err != nil {
		return err
	}

	// wait for session id and voice server
	go func() {
		sessionId, voiceServer := <-sessionIdChan, <-voiceServerChan
		bh_context.GetGatewayClient().CleanVoiceInstantiateChan(guildId)
		client.Lock()
		defer client.Unlock()
		client.sessionId = &sessionId
		client.voiceServer = voiceServer
		go client.connect()
		go client.play()
	}()

	return nil
}

// StopClient remove client from list and properly leave
func (clm *clientManager) StopClient(guildId types.Snowflake) error {
	err := clm.removeClient(guildId)
	if err != nil {
		return err
	}
	err = bh_context.GetGatewayClient().LeaveVoiceChannelMsg(guildId)
	if err != nil {
		return err
	}
	return nil
}
