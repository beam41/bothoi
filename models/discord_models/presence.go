package discord_models

import "bothoi/models/types"

type PresenceUpdate struct {
	User         User            `json:"user" mapstructure:"user"`
	GuildId      types.Snowflake `json:"guild_id" mapstructure:"guild_id"`
	Status       string          `json:"status" mapstructure:"status"`
	Activities   []Activity      `json:"activities" mapstructure:"activities"`
	ClientStatus ClientStatus    `json:"client_status" mapstructure:"client_status"`
}
