package voice

import (
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"bothoi/repo"
	"context"
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

// ClientStart start the client if not started already
func (clm *ClientManager) ClientStart(guildID types.Snowflake) (bool, chan<- string, chan<- *discord_models.VoiceServer) {
	clm.RLock()
	var cli = clm.list[guildID]
	if cli != nil {
		clm.RUnlock()
		cli.RLock()
		defer cli.RUnlock()
		if cli.running {
			go cli.play()
		}
		return false, nil, nil
	}
	clm.RUnlock()

	log.Println(guildID, "Starting client")
	clm.Lock()
	defer clm.Unlock()
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
	sessionIDChan := make(chan string)
	voiceServerChan := make(chan *discord_models.VoiceServer)

	// wait for session id and voice server
	go func() {
		sessionID, voiceServer := <-sessionIDChan, <-voiceServerChan
		clm.list[guildID].Lock()
		defer clm.list[guildID].Unlock()
		clm.list[guildID].sessionID = &sessionID
		clm.list[guildID].voiceServer = voiceServer
		go clm.list[guildID].connect()
		go clm.list[guildID].play()
	}()
	return true, sessionIDChan, voiceServerChan
}

// ClientStop remove client from list and properly leave
func (clm *ClientManager) ClientStop(guildID types.Snowflake) bool {
	_ = repo.DeleteSongsInGuild(guildID)
	clm.Lock()
	defer clm.Unlock()
	var cli = clm.list[guildID]
	if cli == nil {
		return false
	}
	cli.Lock()
	defer cli.Unlock()
	cli.destroyed = true
	cli.udpReadyWait.Broadcast()
	cli.pauseWait.Broadcast()
	cli.ctxCancel()
	cli.connCloseNormal()
	delete(clm.list, guildID)
	return true
}

// ClientPauseSong pause/resume the music player return true if the player is paused
func (clm *ClientManager) ClientPauseSong(guildID types.Snowflake) (found, pausing bool) {
	clm.Lock()
	defer clm.Unlock()
	var cli = clm.list[guildID]
	if cli == nil {
		return false, false
	}
	cli.RLock()
	defer cli.RUnlock()
	cli.pauseWait.L.Lock()
	cli.pausing = !cli.pausing
	if !cli.pausing {
		cli.pauseWait.Broadcast()
	}
	cli.pauseWait.L.Unlock()
	return true, cli.pausing
}

// ClientSkipSong skip a song
func (clm *ClientManager) ClientSkipSong(guildID types.Snowflake) bool {
	clm.Lock()
	defer clm.Unlock()
	var cli = clm.list[guildID]
	if cli == nil {
		return false
	}
	cli.RLock()
	defer cli.RUnlock()
	cli.skip = true
	cli.pauseWait.Broadcast()
	return true
}
