package discord_models

import "bothoi/models/types"

type GuildScheduledEvent struct {
	Id                 types.Snowflake                    `json:"id" mapstructure:"id"`
	GuildId            types.Snowflake                    `json:"guild_id" mapstructure:"guild_id"`
	ChannelId          *types.Snowflake                   `json:"channel_id" mapstructure:"channel_id"`
	CreatorId          *types.Snowflake                   `json:"creator_id" mapstructure:"creator_id"`
	Name               string                             `json:"name" mapstructure:"name"`
	Description        *string                            `json:"description" mapstructure:"description"`
	ScheduledStartTime types.ISOTimeStamp                 `json:"scheduled_start_time" mapstructure:"scheduled_start_time"`
	ScheduledEndTime   *types.ISOTimeStamp                `json:"scheduled_end_time" mapstructure:"scheduled_end_time"`
	PrivacyLevel       uint8                              `json:"privacy_level" mapstructure:"privacy_level"`
	Status             uint8                              `json:"status" mapstructure:"status"`
	EntityType         uint8                              `json:"entity_type" mapstructure:"entity_type"`
	EntityId           *types.Snowflake                   `json:"entity_id" mapstructure:"entity_id"`
	EntityMetadata     *GuildScheduledEventEntityMetadata `json:"entity_metadata" mapstructure:"entity_metadata"`
	Creator            User                               `json:"creator" mapstructure:"creator"`
	UserCount          uint32                             `json:"user_count" mapstructure:"user_count"`
	Image              *string                            `json:"image" mapstructure:"image"`
}

type GuildScheduledEventEntityMetadata struct {
	Location string `json:"location" mapstructure:"location"`
}