package discord_models

import "bothoi/models/types"

type GuildMember struct {
	User                       User                `json:"user" mapstructure:"user"`
	Nick                       *string             `json:"nick" mapstructure:"nick"`
	Avatar                     *string             `json:"avatar" mapstructure:"avatar"`
	Roles                      []types.Snowflake   `json:"roles,string" mapstructure:"roles"`
	JoinedAt                   types.ISOTimeStamp  `json:"joined_at" mapstructure:"joined_at"`
	PremiumSince               *types.ISOTimeStamp `json:"premium_since" mapstructure:"premium_since"`
	Deaf                       bool                `json:"deaf" mapstructure:"deaf"`
	Mute                       bool                `json:"mute" mapstructure:"mute"`
	Pending                    bool                `json:"pending" mapstructure:"pending"`
	Permissions                string              `json:"string" mapstructure:"string"`
	CommunicationDisabledUntil *types.ISOTimeStamp `json:"communication_disabled_until" mapstructure:"communication_disabled_until"`
}
