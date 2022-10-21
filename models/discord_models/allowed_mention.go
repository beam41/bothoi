package discord_models

import "bothoi/models/types"

type AllowedMention struct {
	Parse       []string          `json:"parse,omitempty" mapstructure:"parse"`
	Roles       []types.Snowflake `json:"roles,omitempty" mapstructure:"roles"`
	Users       []types.Snowflake `json:"users,omitempty" mapstructure:"users"`
	RepliedUser bool              `json:"replied_user" mapstructure:"replied_user"`
}
