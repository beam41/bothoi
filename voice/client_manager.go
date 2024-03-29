package voice

import (
	"bothoi/gateway"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"bothoi/repo"
	"context"
	"log"
	"sync"
)

type ClientManager struct {
	sync.RWMutex
	gatewayClient *gateway.Client
	list          map[types.Snowflake]*client
}

func NewClientManager(gatewayClient *gateway.Client) *ClientManager {
	clm := &ClientManager{
		gatewayClient: gatewayClient,
		list:          make(map[types.Snowflake]*client),
	}
	gatewayClient.RegisterNewSessionIDHandler(clm.NewSessionIDHandler)
	return clm
}

func (clm *ClientManager) NewSessionIDHandler(guildID types.Snowflake, sessionID string) {
	clm.Lock()
	defer clm.Unlock()
	var cli = clm.list[guildID]
	if cli == nil {
		return
	}
	cli.Lock()
	defer cli.Unlock()
	cli.sessionID = &sessionID
}

// ClientStart start the client if not started already
func (clm *ClientManager) ClientStart(guildID, channelID types.Snowflake) error {
	clm.RLock()
	var cli = clm.list[guildID]
	if cli != nil {
		clm.RUnlock()
		cli.RLock()
		defer cli.RUnlock()
		if cli.running {
			go cli.play()
		}
		return nil
	}
	clm.RUnlock()

	log.Println(guildID, "Starting client")
	clm.Lock()
	defer clm.Unlock()
	ctx, cancel := context.WithCancel(context.Background())
	clm.list[guildID] = &client{
		guildID:         guildID,
		channelID:       channelID,
		udpReadyWait:    sync.NewCond(&sync.Mutex{}),
		ctx:             ctx,
		ctxCancel:       cancel,
		pauseWait:       sync.NewCond(&sync.Mutex{}),
		stopWaitForExit: make(chan struct{}),
		clm:             clm,
		waitResume:      make(chan struct{}, 1),
	}
	sessionIDChan := make(chan string)
	voiceServerChan := make(chan *discord_models.VoiceServer)
	err := clm.gatewayClient.VoiceChannelJoin(guildID, channelID, sessionIDChan, voiceServerChan)
	if err != nil {
		log.Println(guildID, err)
		return err
	}

	// wait for session id and voice server
	go func() {
		sessionID := <-sessionIDChan
		voiceServer := <-voiceServerChan
		clm.gatewayClient.CleanVoiceInstantiateChan(guildID)
		clm.list[guildID].Lock()
		defer clm.list[guildID].Unlock()
		clm.list[guildID].sessionID = &sessionID
		clm.list[guildID].voiceServer = voiceServer
		go clm.list[guildID].connect()
		go clm.list[guildID].play()
	}()
	return nil
}

// ClientStop remove client from list and properly leave
func (clm *ClientManager) ClientStop(guildID types.Snowflake) error {
	_ = repo.DeleteSongsInGuild(guildID)
	clm.Lock()
	defer clm.Unlock()
	var cli = clm.list[guildID]
	if cli == nil {
		return nil
	}
	cli.Lock()
	defer cli.Unlock()
	cli.destroyed = true
	cli.udpReadyWait.Broadcast()
	cli.pauseWait.Broadcast()
	cli.ctxCancel()
	cli.connCloseNormal()
	delete(clm.list, guildID)
	err := clm.gatewayClient.VoiceChannelLeave(guildID)
	return err
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
