package discord_models

import "bothoi/models/types"

type PresenceUpdate struct {
	User         User            `json:"user" mapstructure:"user"`
	GuildID      types.Snowflake `json:"guild_id,string" mapstructure:"guild_id"`
	Status       string          `json:"status" mapstructure:"status"`
	Activities   []Activity      `json:"activities" mapstructure:"activities"`
	ClientStatus ClientStatus    `json:"client_status" mapstructure:"client_status"`
}
