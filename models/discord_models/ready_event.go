package discord_models

import "bothoi/models/types"

type ReadyEvent struct {
	V           uint8              `json:"v" mapstructure:"v"`
	User        User               `json:"user" mapstructure:"user"`
	Shard       [2]uint16          `json:"shard" mapstructure:"shard"`
	SessionId   types.Snowflake    `json:"session_id,string" mapstructure:"session_id"`
	Guilds      []UnavailableGuild `json:"guilds" mapstructure:"guilds"`
	Application Application        `json:"application" mapstructure:"application"`
}
