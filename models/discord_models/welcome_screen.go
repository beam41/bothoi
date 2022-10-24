package discord_models

import "bothoi/models/types"

type WelcomeScreen struct {
	Description     *string                `json:"description" mapstructure:"description"`
	WelcomeChannels []WelcomeScreenChannel `json:"welcome_channels" mapstructure:"welcome_channels"`
}

type WelcomeScreenChannel struct {
	ChannelId   types.Snowflake  `json:"channel_id,string" mapstructure:"channel_id"`
	Description string           `json:"description" mapstructure:"description"`
	EmojiId     *types.Snowflake `json:"emoji_id,string" mapstructure:"emoji_id"`
	EmojiName   *string          `json:"emoji_name" mapstructure:"emoji_name"`
}
