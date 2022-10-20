package voice

import (
	"bothoi/models"
	"bothoi/util"
	"context"
	"errors"
	"sync"
)

type voiceClientT struct {
	sync.RWMutex
	c map[string]*VoiceClient
}

var clientList = &voiceClientT{
	c: map[string]*VoiceClient{},
}

func addGuildToClient(guildID, voiceChannelId string) *VoiceClient {
	if clientList.c[guildID] != nil {
		return clientList.c[guildID]
	}
	clientList.Lock()
	ctx, cancel := context.WithCancel(context.Background())
	clientList.c[guildID] = &VoiceClient{
		guildID:        guildID,
		voiceChannelID: voiceChannelId,
		songQueue:      []*models.SongItemWData{},
		udpReadyWait:   sync.NewCond(&sync.Mutex{}),
		ctx:            ctx,
		ctxCancel:      cancel,
		frameData:      make(chan []byte),
		pauseWait:      sync.NewCond(&sync.Mutex{}),
	}
	clientList.Unlock()
	return clientList.c[guildID]
}

// add song to the song queue and start playing if not play already
func AppendSongToSongQueue(guildID string, songItem models.SongItem) int {
	clientList.RLock()
	defer clientList.RUnlock()
	var client = clientList.c[guildID]
	if client == nil {
		return 0
	}
	client.Lock()
	defer client.Unlock()
	client.songQueue = append(client.songQueue, &models.SongItemWData{
		YtID:        songItem.YtID,
		Title:       songItem.Title,
		Duration:    songItem.Duration,
		RequesterID: songItem.RequesterID,
	})
	go client.play()
	go client.downloadUpcoming()
	return len(client.songQueue)
}

// stop playing and remove the client from the list
func removeClient(guildID string) error {
	clientList.Lock()
	defer clientList.Unlock()
	var client = clientList.c[guildID]
	if client == nil {
		return errors.New("Client not found")
	}
	client.Lock()
	defer client.Unlock()
	client.destroyed = true
	client.udpReadyWait.Broadcast()
	client.pauseWait.Broadcast()
	client.ctxCancel()
	delete(clientList.c, guildID)
	close(client.frameData)
	return nil
}

// get the copy of current song queue
func GetSongQueue(guildID string, start, end int) (playing bool, queue []models.SongItem) {
	clientList.RLock()
	defer clientList.RUnlock()
	var client = clientList.c[guildID]
	if client == nil {
		return false, nil
	}
	client.RLock()
	defer client.RUnlock()
	queue = make([]models.SongItem, end-start)
	for i, item := range client.songQueue[start:util.Min(len(client.songQueue), end)] {
		queue[i] = models.SongItem{
			YtID:        item.YtID,
			Title:       item.Title,
			Duration:    item.Duration,
			RequesterID: item.RequesterID,
		}
	}
	return client.playing && !client.pausing, queue
}

func ClientExist(guildID string) bool {
	clientList.RLock()
	defer clientList.RUnlock()
	var client = clientList.c[guildID]
	return client != nil
}

// pause/resume the music player return true if the player is paused
func PauseClient(guildID string) (bool, error) {
	clientList.Lock()
	defer clientList.Unlock()
	var client = clientList.c[guildID]
	if client == nil {
		return false, errors.New("Client not found")
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
func SkipSong(guildID string) error {
	clientList.Lock()
	defer clientList.Unlock()
	var client = clientList.c[guildID]
	if client == nil {
		return errors.New("Client not found")
	}
	client.RLock()
	defer client.RUnlock()
	client.skip = true
	client.pauseWait.Broadcast()
	return nil
}

func GetVoiceChannelID(guildID string) string {
	clientList.RLock()
	defer clientList.RUnlock()
	var client = clientList.c[guildID]
	if client == nil {
		return ""
	}
	client.RLock()
	defer client.RUnlock()
	return client.voiceChannelID
}
