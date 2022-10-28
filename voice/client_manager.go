package voice

import (
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"bothoi/repo"
	"context"
	"errors"
	"log"
	"sync"
)

type ClientManager struct {
	sync.RWMutex
	list map[types.Snowflake]*client
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		list: make(map[types.Snowflake]*client),
	}
}

// createClient make new client from guild
func (clm *ClientManager) createClient(guildID types.Snowflake) *client {
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
func (clm *ClientManager) removeClient(guildID types.Snowflake) error {
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
func (clm *ClientManager) PauseClient(guildID types.Snowflake) (bool, error) {
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
func (clm *ClientManager) SkipSong(guildID types.Snowflake) error {
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
func (clm *ClientManager) StartClient(guildID types.Snowflake) (bool, chan<- string, chan<- *discord_models.VoiceServer) {
	clm.RLock()
	var client = clm.list[guildID]
	if client != nil {
		clm.RUnlock()
		client.RLock()
		defer client.RUnlock()
		if client.running {
			go client.play()
		}
		return false, nil, nil
	}
	clm.RUnlock()

	log.Println(guildID, "Starting client")
	client = clm.createClient(guildID)
	sessionIDChan := make(chan string)
	voiceServerChan := make(chan *discord_models.VoiceServer)

	// wait for session id and voice server
	go func() {
		sessionID, voiceServer := <-sessionIDChan, <-voiceServerChan
		client.Lock()
		defer client.Unlock()
		client.sessionID = &sessionID
		client.voiceServer = voiceServer
		go client.connect()
		go client.play()
	}()
	return true, sessionIDChan, voiceServerChan
}

// StopClient remove client from list and properly leave
func (clm *ClientManager) StopClient(guildID types.Snowflake) error {
	_ = repo.DeleteSongsInGuild(guildID)
	err := clm.removeClient(guildID)
	if err != nil {
		return err
	}
	return nil
}
