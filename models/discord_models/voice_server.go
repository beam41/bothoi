package discord_models

import "bothoi/models/types"

type VoiceServer struct {
	Token    string          `json:"token" mapstructure:"token"`
	GuildID  types.Snowflake `json:"guild_id,string" mapstructure:"guild_id"`
	Endpoint string          `json:"endpoint" mapstructure:"endpoint"`
}
