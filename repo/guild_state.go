package repo

import (
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"sync"
)

type guildStateT struct {
	sync.RWMutex
	state map[types.Snowflake]*discord_models.GuildCreate
}

var guildState = &guildStateT{
	state: map[types.Snowflake]*discord_models.GuildCreate{},
}

func AddGuild(guild *discord_models.GuildCreate) {
	guildState.Lock()
	guildState.state[guild.Id] = guild
	guildState.Unlock()
}
