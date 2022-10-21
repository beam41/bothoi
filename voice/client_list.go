package voice

import (
	"bothoi/models"
	"bothoi/models/types"
	"bothoi/util"
	"context"
	"errors"
	"sync"
)

type voiceClientT struct {
	sync.RWMutex
	c map[types.Snowflake]*VClient
}

var clientList = &voiceClientT{
	c: map[types.Snowflake]*VClient{},
}

func addGuildToClient(guildId, voiceChannelId types.Snowflake) *VClient {
	if clientList.c[guildId] != nil {
		return clientList.c[guildId]
	}
	clientList.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	clientList.c[guildId] = &VClient{
		guildId:         guildId,
		voiceChannelId:  voiceChannelId,
		songQueue:       []*models.SongItemWData{},
		udpReadyWait:    sync.NewCond(&sync.Mutex{}),
		ctx:             ctx,
		ctxCancel:       cancel,
		frameData:       make(chan []byte),
		pauseWait:       sync.NewCond(&sync.Mutex{}),
		stopWaitForExit: make(chan struct{}),
	}
	clientList.Unlock()
	return clientList.c[guildId]
}

// add song to the song queue and start playing if not play already
func AppendSongToSongQueue(guildId types.Snowflake, songItem models.SongItem) int {
	clientList.RLock()
	defer clientList.RUnlock()
	var client = clientList.c[guildId]
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

// stop playing and remove the client from the list
func removeClient(guildId types.Snowflake) error {
	clientList.Lock()
	defer clientList.Unlock()
	var client = clientList.c[guildId]
	if client == nil {
		return errors.New("client not found")
	}
	client.Lock()
	defer client.Unlock()
	client.destroyed = true
	client.udpReadyWait.Broadcast()
	client.pauseWait.Broadcast()
	client.ctxCancel()
	delete(clientList.c, guildId)
	close(client.frameData)
	return nil
}

// get the copy of current song queue
func GetSongQueue(guildId types.Snowflake, start, end int) (playing bool, queue []models.SongItem) {
	clientList.RLock()
	defer clientList.RUnlock()
	var client = clientList.c[guildId]
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

func ClientExist(guildId types.Snowflake) bool {
	clientList.RLock()
	defer clientList.RUnlock()
	var client = clientList.c[guildId]
	return client != nil
}

// pause/resume the music player return true if the player is paused
func PauseClient(guildId types.Snowflake) (bool, error) {
	clientList.Lock()
	defer clientList.Unlock()
	var client = clientList.c[guildId]
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

// skip a song
func SkipSong(guildId types.Snowflake) error {
	clientList.Lock()
	defer clientList.Unlock()
	var client = clientList.c[guildId]
	if client == nil {
		return errors.New("client not found")
	}
	client.RLock()
	defer client.RUnlock()
	client.skip = true
	client.pauseWait.Broadcast()
	return nil
}

func GetVoiceChannelId(guildId types.Snowflake) types.Snowflake {
	clientList.RLock()
	defer clientList.RUnlock()
	var client = clientList.c[guildId]
	if client == nil {
		return ""
	}
	client.RLock()
	defer client.RUnlock()
	return client.voiceChannelId
}
