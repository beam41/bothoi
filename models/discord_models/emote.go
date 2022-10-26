package discord_models

import "bothoi/models/types"

type Emoji struct {
	ID            types.Snowflake `json:"id,string" mapstructure:"id"`
	Name          *string         `json:"name" mapstructure:"name"`
	Roles         []string        `json:"roles" mapstructure:"roles"`
	User          User            `json:"user" mapstructure:"user"`
	RequireColons bool            `json:"require_colons" mapstructure:"require_colons"`
	Managed       bool            `json:"managed" mapstructure:"managed"`
	Animated      bool            `json:"animated" mapstructure:"animated"`
	Available     bool            `json:"available" mapstructure:"available"`
}

type Sticker struct {
	ID          string          `json:"id" mapstructure:"id"`
	PackID      types.Snowflake `json:"pack_id,string" mapstructure:"pack_id"`
	Name        string          `json:"name" mapstructure:"name"`
	Description *string         `json:"description" mapstructure:"description"`
	Tags        *string         `json:"tags" mapstructure:"tags"`
	Type        uint8           `json:"type" mapstructure:"type"`
	FormatType  uint8           `json:"format_type" mapstructure:"format_type"`
	Available   bool            `json:"available" mapstructure:"available"`
	GuildID     types.Snowflake `json:"guild_id,string" mapstructure:"guild_id"`
	User        User            `json:"user" mapstructure:"user"`
	SortValue   uint16          `json:"sort_value" mapstructure:"sort_value"`
}
