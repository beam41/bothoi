package discord_models

import "bothoi/models/types"

type Emoji struct {
	Id            types.Snowflake `json:"id" mapstructure:"id"`
	Name          *string         `json:"name" mapstructure:"name"`
	Roles         []string        `json:"roles" mapstructure:"roles"`
	User          User            `json:"user" mapstructure:"user"`
	RequireColons bool            `json:"require_colons" mapstructure:"require_colons"`
	Managed       bool            `json:"managed" mapstructure:"managed"`
	Animated      bool            `json:"animated" mapstructure:"animated"`
	Available     bool            `json:"available" mapstructure:"available"`
}

type Sticker struct {
	Id          string          `json:"id" mapstructure:"id"`
	PackId      types.Snowflake `json:"pack_id" mapstructure:"pack_id"`
	Name        string          `json:"name" mapstructure:"name"`
	Description *string         `json:"description" mapstructure:"description"`
	Tags        *string         `json:"tags" mapstructure:"tags"`
	Type        uint8           `json:"type" mapstructure:"type"`
	FormatType  uint8           `json:"format_type" mapstructure:"format_type"`
	Available   bool            `json:"available" mapstructure:"available"`
	GuildId     types.Snowflake `json:"guild_id" mapstructure:"guild_id"`
	User        User            `json:"user" mapstructure:"user"`
	SortValue   uint16          `json:"sort_value" mapstructure:"sort_value"`
}
