package voice

import (
	"bothoi/models"
	"bothoi/states"
	"sync"
)

type voiceChanMapChan struct {
	sessionIdChan   chan<- string
	voiceServerChan chan<- *models.VoiceServer
}

var voiceChanMap = map[string]voiceChanMapChan{}
var voiceChanMapMutex sync.RWMutex

func ReturnSessionId(guildID, sessionID string) {
	voiceChanMapMutex.RLock()
	defer voiceChanMapMutex.RUnlock()
	if chanMap, ok := voiceChanMap[guildID]; ok {
		chanMap.sessionIdChan <- sessionID
	}
}

func ReturnVoiceServer(guildID string, voiceServer *models.VoiceServer) {
	voiceChanMapMutex.RLock()
	defer voiceChanMapMutex.RUnlock()
	if chanMap, ok := voiceChanMap[guildID]; ok {
		chanMap.voiceServerChan <- voiceServer
	}
}

func joinVoiceChannel(guildID, channelID string, sessionIdChan chan<- string, voiceServerChan chan<- *models.VoiceServer) error {
	createVoice := models.NewVoiceStateUpdate(guildID, &channelID, false, true)
	err := states.GatewayConnWriteJSON(createVoice)
	if err != nil {
		return err
	}
	voiceChanMapMutex.Lock()
	defer voiceChanMapMutex.Unlock()
	voiceChanMap[guildID] = voiceChanMapChan{sessionIdChan, voiceServerChan}
	return nil
}

func leaveVoiceChannel(guildID string) error {
	leaveVoice := models.NewVoiceStateUpdate(guildID, nil, false, false)
	err := states.GatewayConnWriteJSON(leaveVoice)
	if err != nil {
		return err
	}
	voiceChanMapMutex.Lock()
	defer voiceChanMapMutex.Unlock()
	delete(voiceChanMap, guildID)
	return nil
}
