package discord_models

import "bothoi/models/types"

type StageInstance struct {
	ID                    types.Snowflake  `json:"id,string" mapstructure:"id"`
	GuildID               types.Snowflake  `json:"guild_id,string" mapstructure:"guild_id"`
	ChannelID             types.Snowflake  `json:"channel_id,string" mapstructure:"channel_id"`
	Topic                 string           `json:"topic" mapstructure:"topic"`
	PrivacyLevel          uint8            `json:"privacy_level" mapstructure:"privacy_level"`
	DiscoverableDisabled  bool             `json:"discoverable_disabled" mapstructure:"discoverable_disabled"`
	GuildScheduledEventID *types.Snowflake `json:"guild_scheduled_event_id,string" mapstructure:"guild_scheduled_event_id"`
}
