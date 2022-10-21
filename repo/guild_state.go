package repo

import (
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"sync"
)

type guildStateT struct {
	sync.RWMutex
	state map[types.Snowflake]*discord_models.Guild
}

var guildState = &guildStateT{
	state: map[types.Snowflake]*discord_models.Guild{},
}

func AddGuild(guild *discord_models.Guild) {
	guildState.Lock()
	guildState.state[guild.Id] = guild
	guildState.Unlock()
}
