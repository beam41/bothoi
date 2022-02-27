package states

import (
	"bothoi/models"
	"sync"
)

// map[GuildID]Guild
var GuildState = map[string]*models.Guild{}

var GuildStateLock sync.Mutex

func AddGuild(guild *models.Guild) {
	GuildStateLock.Lock()
	defer GuildStateLock.Unlock()
	GuildState[guild.ID] = guild
}
