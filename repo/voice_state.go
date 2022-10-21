package repo

import (
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"sync"
)

type voiceStateT struct {
	sync.RWMutex
	state map[types.Snowflake]*discord_models.VoiceState
}

var voiceState_ = &voiceStateT{
	state: map[types.Snowflake]*discord_models.VoiceState{},
}

func AddVoiceState(voiceState *discord_models.VoiceState) {
	voiceState_.Lock()
	voiceState_.state[voiceState.UserId] = voiceState
	voiceState_.Unlock()
}

func AddVoiceStateBulk(voiceStates []discord_models.VoiceState) {
	voiceState_.Lock()
	for _, voiceState := range voiceStates {
		voiceState_.state[voiceState.UserId] = &voiceState
	}
	voiceState_.Unlock()
}

func GetVoiceState(userId types.Snowflake) *discord_models.VoiceState {
	voiceState_.RLock()
	defer voiceState_.RUnlock()
	return voiceState_.state[userId]
}
