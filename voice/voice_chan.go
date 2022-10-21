package voice

import (
	"bothoi/global"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"sync"
)

type voiceChanMapChan struct {
	sessionIdChan   chan<- string
	voiceServerChan chan<- *discord_models.VoiceServer
}

var voiceChanMap = map[types.Snowflake]voiceChanMapChan{}
var voiceChanMapMutex sync.RWMutex

func ReturnSessionId(guildId types.Snowflake, sessionId string) {
	voiceChanMapMutex.RLock()
	defer voiceChanMapMutex.RUnlock()
	if chanMap, ok := voiceChanMap[guildId]; ok {
		chanMap.sessionIdChan <- sessionId
	}
}

func ReturnVoiceServer(guildId types.Snowflake, voiceServer *discord_models.VoiceServer) {
	voiceChanMapMutex.RLock()
	defer voiceChanMapMutex.RUnlock()
	if chanMap, ok := voiceChanMap[guildId]; ok {
		chanMap.voiceServerChan <- voiceServer
	}
}

func joinVoiceChannel(guildId, channelId types.Snowflake, sessionIdChan chan<- string, voiceServerChan chan<- *discord_models.VoiceServer) error {
	createVoice := discord_models.NewVoiceStateUpdate(guildId, &channelId, false, true)
	err := global.GatewayConnWriteJSON(createVoice)
	if err != nil {
		return err
	}
	voiceChanMapMutex.Lock()
	defer voiceChanMapMutex.Unlock()
	voiceChanMap[guildId] = voiceChanMapChan{sessionIdChan, voiceServerChan}
	return nil
}

func leaveVoiceChannel(guildId types.Snowflake) error {
	leaveVoice := discord_models.NewVoiceStateUpdate(guildId, nil, false, false)
	err := global.GatewayConnWriteJSON(leaveVoice)
	if err != nil {
		return err
	}
	voiceChanMapMutex.Lock()
	defer voiceChanMapMutex.Unlock()
	delete(voiceChanMap, guildId)
	return nil
}
