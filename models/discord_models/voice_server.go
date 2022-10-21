package discord_models

import "bothoi/models/types"

type VoiceServer struct {
	Token    string          `json:"token" mapstructure:"token"`
	GuildId  types.Snowflake `json:"guild_id" mapstructure:"guild_id"`
	Endpoint string          `json:"endpoint" mapstructure:"endpoint"`
}
