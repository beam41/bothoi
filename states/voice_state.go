package states

import (
	"bothoi/models"
	"sync"
)

// map[UserID]voicestate
var VoiceState = map[string]*models.VoiceState{}

var VoiceStateLock sync.Mutex

func AddVoiceState(voiceState *models.VoiceState) {
	VoiceStateLock.Lock()
	defer VoiceStateLock.Unlock()
	VoiceState[voiceState.UserID] = voiceState
}

func AddVoiceStateBulk(voiceStates []models.VoiceState) {
	VoiceStateLock.Lock()
	defer VoiceStateLock.Unlock()
	for _, voiceState := range voiceStates {
		VoiceState[voiceState.UserID] = &voiceState
	}
}
