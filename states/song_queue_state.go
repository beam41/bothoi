package states

import (
	"bothoi/models"
	"sync"
)

type songQueueStateT struct {
	sync.RWMutex
	state map[string]*models.SongQueue
}

var songQueueState = &songQueueStateT{
	state: map[string]*models.SongQueue{},
}

func AddGuildToSongQueue(songQueue *models.SongQueue) {
	songQueueState.Lock()
	songQueueState.state[songQueue.GuildID] = songQueue
	songQueueState.Unlock()
}

func AppendSongToSongQueue(guildID string, songItem models.SongItem) {
	songQueueState.Lock()
	var songQ = songQueueState.state[guildID]
	_ = append(songQ.Songs, songItem)
	songQueueState.Unlock()
}

func RemoveGuildFromSongQueue(guildID string) {
	songQueueState.Lock()
	delete(songQueueState.state, guildID)
	songQueueState.Unlock()
}

func SetSessionId(guildID, sessionID string) {
	songQueueState.Lock()
	var songQ = songQueueState.state[guildID]
	songQ.SessionID = &sessionID
	songQueueState.Unlock()
}

func SetVoiceServer(guildID string, voiceServer *models.VoiceServer) {
	songQueueState.Lock()
	var songQ = songQueueState.state[guildID]
	songQ.VoiceServer = voiceServer
	songQueueState.Unlock()
}

func GetSongQueue(guildID string) *models.SongQueue {
	songQueueState.RLock()
	defer songQueueState.RUnlock()
	return songQueueState.state[guildID]
}
