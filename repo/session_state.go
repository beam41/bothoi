package repo

import (
	"bothoi/models/discord_models"
	"sync"
)

var sessionState struct {
	sync.RWMutex
	state *discord_models.ReadyEvent
}

func AddSessionState(state *discord_models.ReadyEvent) {
	sessionState.Lock()
	sessionState.state = state
	sessionState.Unlock()
}

func GetSessionState() *discord_models.ReadyEvent {
	sessionState.RLock()
	defer sessionState.RUnlock()
	return sessionState.state
}
