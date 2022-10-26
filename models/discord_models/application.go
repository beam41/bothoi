package discord_models

import "bothoi/models/types"

type Application struct {
	ID    types.Snowflake `json:"id,string" mapstructure:"id"`
	Flags uint32          `json:"flags" mapstructure:"flags"`
}
