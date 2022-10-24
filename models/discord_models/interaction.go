package discord_models

import "bothoi/models/types"

type Interaction struct {
	Id            types.Snowflake `json:"id,string" mapstructure:"id"`
	ApplicationId types.Snowflake `json:"application_id,string" mapstructure:"application_id"`
	Type          uint8           `json:"type" mapstructure:"type"`
	Data          InteractionData `json:"data" mapstructure:"data"`
	GuildId       types.Snowflake `json:"guild_id,string" mapstructure:"guild_id"`
	ChannelId     types.Snowflake `json:"channel_id,string" mapstructure:"channel_id"`
	Member        GuildMember     `json:"member" mapstructure:"member"`
	Token         string          `json:"token" mapstructure:"token"`
	Version       uint8           `json:"version" mapstructure:"version"`
	Locale        string          `json:"locale" mapstructure:"locale"`
	GuildLocale   string          `json:"guild_locale" mapstructure:"guild_locale"`
}

type InteractionData struct {
	Id      types.Snowflake     `json:"id,string" mapstructure:"id"`
	Name    string              `json:"name" mapstructure:"name"`
	Options []InteractionOption `json:"options" mapstructure:"options"`
	Type    uint8               `json:"type" mapstructure:"type"`
}

type InteractionOption struct {
	Name  string `json:"name" mapstructure:"name"`
	Type  uint8  `json:"type" mapstructure:"type"`
	Value any    `json:"value" mapstructure:"value"`
}
