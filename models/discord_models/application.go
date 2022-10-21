package discord_models

import "bothoi/models/types"

type Application struct {
	Id    types.Snowflake `json:"id" mapstructure:"id"`
	Flags uint32          `json:"flags" mapstructure:"flags"`
}
