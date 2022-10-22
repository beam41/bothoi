package voice

import (
	"bothoi/bh_context"
	"bothoi/models"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"bothoi/util"
	"bothoi/voice/voice_interface"
	"context"
	"errors"
	"sync"
)

type clientManager struct {
	listLock sync.RWMutex
	list     map[types.Snowflake]*client
}

func NewClientManager() voice_interface.ClientManagerInterface {
	return &clientManager{
		list: make(map[types.Snowflake]*client),
	}
}

// createClient make new client from guild
func (clm *clientManager) createClient(guildId, voiceChannelId types.Snowflake) *client {
	if clm.list[guildId] != nil {
		return clm.list[guildId]
	}
	clm.listLock.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	clm.list[guildId] = &client{
		guildId:         guildId,
		voiceChannelId:  voiceChannelId,
		songQueue:       []*models.SongItemWData{},
		udpReadyWait:    sync.NewCond(&sync.Mutex{}),
		ctx:             ctx,
		ctxCancel:       cancel,
		frameData:       make(chan []byte),
		pauseWait:       sync.NewCond(&sync.Mutex{}),
		stopWaitForExit: make(chan struct{}),
		clm:             clm,
	}
	clm.listLock.Unlock()
	return clm.list[guildId]
}

// AppendSongToSongQueue add song to the song queue and start playing if not play already
func (clm *clientManager) AppendSongToSongQueue(guildId types.Snowflake, songItem models.SongItem) int {
	clm.listLock.RLock()
	defer clm.listLock.RUnlock()
	var client = clm.list[guildId]
	if client == nil {
		return 0
	}
	client.Lock()
	defer client.Unlock()
	client.songQueue = append(client.songQueue, &models.SongItemWData{
		YtId:        songItem.YtId,
		Title:       songItem.Title,
		Duration:    songItem.Duration,
		RequesterId: songItem.RequesterId,
	})
	go client.play()
	go client.downloadUpcoming()
	return len(client.songQueue)
}

// removeClient stop playing and remove the client from the list
func (clm *clientManager) removeClient(guildId types.Snowflake) error {
	clm.listLock.Lock()
	defer clm.listLock.Unlock()
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
	close(client.frameData)
	return nil
}

// GetSongQueue get the copy of current song queue
func (clm *clientManager) GetSongQueue(guildId types.Snowflake, start, end int) (playing bool, queue []models.SongItem) {
	clm.listLock.RLock()
	defer clm.listLock.RUnlock()
	var client = clm.list[guildId]
	if client == nil {
		return false, nil
	}
	client.RLock()
	defer client.RUnlock()
	queue = make([]models.SongItem, end-start)
	for i, item := range client.songQueue[start:util.Min(len(client.songQueue), end)] {
		queue[i] = models.SongItem{
			YtId:        item.YtId,
			Title:       item.Title,
			Duration:    item.Duration,
			RequesterId: item.RequesterId,
		}
	}
	return client.playing && !client.pausing, queue
}

// PauseClient pause/resume the music player return true if the player is paused
func (clm *clientManager) PauseClient(guildId types.Snowflake) (bool, error) {
	clm.listLock.Lock()
	defer clm.listLock.Unlock()
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
	clm.listLock.Lock()
	defer clm.listLock.Unlock()
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

func (clm *clientManager) GetVoiceChannelId(guildId types.Snowflake) types.Snowflake {
	clm.listLock.RLock()
	defer clm.listLock.RUnlock()
	var client = clm.list[guildId]
	if client == nil {
		return ""
	}
	client.RLock()
	defer client.RUnlock()
	return client.voiceChannelId
}

// StartClient start the client if not started already
func (clm *clientManager) StartClient(guildId, channelId types.Snowflake) error {
	clm.listLock.RLock()
	var client = clm.list[guildId]
	if client != nil {
		clm.listLock.RUnlock()
		client.RLock()
		defer client.RUnlock()
		if (client.voiceChannelId != "") && (client.voiceChannelId != channelId) {
			return errors.New("already in a different voice channel")
		} else if client.running {
			return nil
		}
	}
	clm.listLock.RUnlock()

	client = clm.createClient(guildId, channelId)
	sessionIdChan := make(chan string)
	voiceServerChan := make(chan *discord_models.VoiceServer)
	err := bh_context.Ctx.GatewayClient.JoinVoiceChannelMsg(guildId, channelId, sessionIdChan, voiceServerChan)
	if err != nil {
		return err
	}

	// wait for session id and voice server
	go func() {
		sessionId, voiceServer := <-sessionIdChan, <-voiceServerChan
		bh_context.Ctx.GatewayClient.CleanVoiceInstantiateChan(guildId)
		client.Lock()
		defer client.Unlock()
		client.sessionId = &sessionId
		client.voiceServer = voiceServer
		go client.connect()
	}()

	return nil
}

// StopClient remove client from list and properly leave
func (clm *clientManager) StopClient(guildId types.Snowflake) error {
	err := clm.removeClient(guildId)
	if err != nil {
		return err
	}
	err = bh_context.Ctx.GatewayClient.LeaveVoiceChannelMsg(guildId)
	if err != nil {
		return err
	}
	return nil
}
