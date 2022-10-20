package gateway

import (
	"bothoi/config"
	"bothoi/models"
	"bothoi/states"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/mitchellh/mapstructure"
)

var executorList = map[string]func(*models.Interaction){}

func mapInteractionExecute(data *models.Interaction) {
	if interaction, ok := executorList[data.Data.Name]; ok {
		interaction(data)
	}
}

func SetExecutorList(list map[string]func(*models.Interaction)) {
	executorList = list
}

func dispatchHandler(c *websocket.Conn, payload models.GatewayPayload) {
	switch payload.T {
	case "READY":
		var sessionState models.SessionState
		mapstructure.Decode(payload.D, &sessionState)
		states.AddSessionState(&sessionState)
	case "INTERACTION_CREATE":
		var data models.Interaction
		mapstructure.Decode(payload.D, &data)
		mapInteractionExecute(&data)
	case "GUILD_CREATE":
		var data models.Guild
		mapstructure.Decode(payload.D, &data)
		// guild voice state don't contain guild id
		states.AddGuild(&data)
		var voiceStates []models.VoiceState
		for _, voiceState := range data.VoiceStates {
			voiceState.GuildID = data.ID
			voiceStates = append(voiceStates, voiceState)
		}
		states.AddVoiceStateBulk(voiceStates)
	case "VOICE_STATE_UPDATE":
		var data *models.VoiceState = new(models.VoiceState)
		mapstructure.Decode(payload.D, data)
		if data.UserID != config.BOT_ID {
			states.AddVoiceState(data)
		} else {
			returnSessionId(data.GuildID, data.SessionID)
		}
	case "VOICE_SERVER_UPDATE":
		var data models.VoiceServer
		mapstructure.Decode(payload.D, &data)
		returnVoiceServer(data.GuildID, &data)
	case "GUILD_UPDATE":
		// not important now
	case "GUILD_DELETE":
		// not important now
	}
}

type voiceChanMapChan struct {
	sessionIdChan   chan<- string
	voiceServerChan chan<- *models.VoiceServer
}

var voiceChanMap map[string]voiceChanMapChan = map[string]voiceChanMapChan{}
var voiceChanMapMutex sync.RWMutex

func returnSessionId(guildID, sessionID string) {
	voiceChanMapMutex.RLock()
	defer voiceChanMapMutex.RUnlock()
	if chanMap, ok := voiceChanMap[guildID]; ok {
		chanMap.sessionIdChan <- sessionID
	}
}

func returnVoiceServer(guildID string, voiceServer *models.VoiceServer) {
	voiceChanMapMutex.RLock()
	defer voiceChanMapMutex.RUnlock()
	if chanMap, ok := voiceChanMap[guildID]; ok {
		chanMap.voiceServerChan <- voiceServer
	}
}
