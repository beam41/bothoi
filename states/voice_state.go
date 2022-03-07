package states

import (
	"bothoi/models"
	"sync"
)

type voiceStateT struct {
	sync.RWMutex
	state  map[string]*models.VoiceState
}

var voiceState_ = &voiceStateT{
	state: map[string]*models.VoiceState{},
}

func AddVoiceState(voiceState *models.VoiceState) {
	voiceState_.Lock()
	voiceState_.state[voiceState.UserID] = voiceState
	voiceState_.Unlock()
}

func AddVoiceStateBulk(voiceStates []models.VoiceState) {
	voiceState_.Lock()
	for _, voiceState := range voiceStates {
		voiceState_.state[voiceState.UserID] = &voiceState
	}
	voiceState_.Unlock()
}

func GetVoiceState(userID string) *models.VoiceState {
	voiceState_.RLock()
	defer voiceState_.RUnlock()
	return voiceState_.state[userID]
}
