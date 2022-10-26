package discord_models

import "bothoi/models/types"

type Role struct {
	ID           types.Snowflake `json:"id,string" mapstructure:"id"`
	Name         string          `json:"name" mapstructure:"name"`
	Color        uint32          `json:"color" mapstructure:"color"`
	Hoist        bool            `json:"hoist" mapstructure:"hoist"`
	Icon         *string         `json:"icon" mapstructure:"icon"`
	UnicodeEmoji *string         `json:"unicode_emoji" mapstructure:"unicode_emoji"`
	Position     uint16          `json:"position" mapstructure:"position"`
	Permissions  string          `json:"permissions" mapstructure:"permissions"`
	Managed      bool            `json:"managed" mapstructure:"managed"`
	Mentionable  bool            `json:"mentionable" mapstructure:"mentionable"`
	Tags         *RoleTags       `json:"tags,omitempty" mapstructure:"tags"`
}

type RoleTags struct {
	BotID             *types.Snowflake `json:"bot_id,string" mapstructure:"bot_id"`
	IntegrationID     *types.Snowflake `json:"integration_id,string" mapstructure:"integration_id"`
	PremiumSubscriber bool             `json:"premium_subscriber" mapstructure:"premium_subscriber"`
}
