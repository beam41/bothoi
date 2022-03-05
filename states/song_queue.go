package states

import (
	"bothoi/models"
	"sync"
)

// map[GuildID]SongQueue
var SongQueue map[string]*models.SongQueue = map[string]*models.SongQueue{}

var SongQueueLock sync.Mutex

func AddGuildToSongQueue(songQueue *models.SongQueue) {
	SongQueueLock.Lock()
	defer SongQueueLock.Unlock()
	SongQueue[songQueue.GuildID] = songQueue
}

func AppendSongToSongQueue(guildID string, songItem models.SongItem) {
	SongQueueLock.Lock()
	defer SongQueueLock.Unlock()
	var songQ = SongQueue[guildID]
	_ = append(songQ.Songs, songItem)
}

func RemoveGuildFromSongQueue(guildID string) {
	SongQueueLock.Lock()
	defer SongQueueLock.Unlock()
	delete(SongQueue, guildID)
}

func SetSessionId(guildID, sessionID string) {
	SongQueueLock.Lock()
	defer SongQueueLock.Unlock()
	var songQ = SongQueue[guildID]
	songQ.SessionID = &sessionID
}

func SetVoiceServer(guildID string, voiceServer *models.VoiceServer) {
	SongQueueLock.Lock()
	defer SongQueueLock.Unlock()
	var songQ = SongQueue[guildID]
	songQ.VoiceServer = voiceServer
}
