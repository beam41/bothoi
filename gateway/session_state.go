package gateway

import (
	"bothoi/models"
	"sync"
)

var ActiveSessionState models.SessionState = models.SessionState{}

var SessionReady sync.WaitGroup;

