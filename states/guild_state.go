package states

import (
	"bothoi/models"
	"sync"
)

type guildStateT struct {
	sync.RWMutex
	state map[string]*models.Guild
}

var guildState = &guildStateT{
	state: map[string]*models.Guild{},
}

func AddGuild(guild *models.Guild) {
	guildState.Lock()
	guildState.state[guild.ID] = guild
	guildState.Unlock()
}
