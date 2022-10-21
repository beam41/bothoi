package discord_models

import "bothoi/models/types"

type Interaction struct {
	Id            types.Snowflake `json:"id" mapstructure:"id"`
	ApplicationId types.Snowflake `json:"application_id" mapstructure:"application_id"`
	Type          uint8           `json:"types" mapstructure:"types"`
	Data          InteractionData `json:"data" mapstructure:"data"`
	GuildId       types.Snowflake `json:"guild_id" mapstructure:"guild_id"`
	ChannelId     types.Snowflake `json:"channel_id" mapstructure:"channel_id"`
	Member        GuildMember     `json:"member" mapstructure:"member"`
	Token         string          `json:"token" mapstructure:"token"`
	Version       uint8           `json:"version" mapstructure:"version"`
	Locale        string          `json:"locale" mapstructure:"locale"`
	GuildLocale   string          `json:"guild_locale" mapstructure:"guild_locale"`
}

type InteractionData struct {
	Id      types.Snowflake     `json:"id" mapstructure:"id"`
	Name    string              `json:"name" mapstructure:"name"`
	Options []InteractionOption `json:"options" mapstructure:"options"`
	Type    uint8               `json:"types" mapstructure:"types"`
}

type InteractionOption struct {
	Name  string `json:"name" mapstructure:"name"`
	Type  uint8  `json:"types" mapstructure:"types"`
	Value any    `json:"value" mapstructure:"value"`
}
