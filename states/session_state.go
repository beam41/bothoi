package states

import (
	"bothoi/models"
	"sync"
)

var SessionState = &models.SessionState{}

var SessionStateReady sync.WaitGroup
