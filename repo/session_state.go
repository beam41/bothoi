package repo

import (
	"bothoi/models"
	"sync"
)

var sessionState struct {
	sync.RWMutex
	state *models.SessionState
}

func AddSessionState(state *models.SessionState) {
	sessionState.Lock()
	sessionState.state = state
	sessionState.Unlock()
}

func GetSessionState() *models.SessionState {
	sessionState.RLock()
	defer sessionState.RUnlock()
	return sessionState.state
}
