package gateway

import (
	"bothoi/models"
	"sync"
)

var activeSessionState models.SessionState = models.SessionState{}

var sessionReady sync.WaitGroup
