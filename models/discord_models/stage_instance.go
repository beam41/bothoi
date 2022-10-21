package discord_models

import "bothoi/models/types"

type StageInstance struct {
	Id                    types.Snowflake  `json:"id" mapstructure:"id"`
	GuildId               types.Snowflake  `json:"guild_id" mapstructure:"guild_id"`
	ChannelId             types.Snowflake  `json:"channel_id" mapstructure:"channel_id"`
	Topic                 string           `json:"topic" mapstructure:"topic"`
	PrivacyLevel          uint8            `json:"privacy_level" mapstructure:"privacy_level"`
	DiscoverableDisabled  bool             `json:"discoverable_disabled" mapstructure:"discoverable_disabled"`
	GuildScheduledEventId *types.Snowflake `json:"guild_scheduled_event_id" mapstructure:"guild_scheduled_event_id"`
}
