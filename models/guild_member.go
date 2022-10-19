package models

type GuildMember struct {
	Deaf                       bool     `json:"deaf" mapstructure:"deaf"`
	HoistedRole                any      `json:"hoisted_role" mapstructure:"hoisted_role"`
	JoinedAt                   string   `json:"joined_at" mapstructure:"joined_at"`
	Mute                       bool     `json:"mute" mapstructure:"mute"`
	Roles                      []string `json:"roles" mapstructure:"roles"`
	User                       User     `json:"user" mapstructure:"user"`
	Avatar                     any      `json:"avatar" mapstructure:"avatar"`
	CommunicationDisabledUntil any      `json:"communication_disabled_until" mapstructure:"communication_disabled_until"`
	Nick                       any      `json:"nick" mapstructure:"nick"`
	Pending                    *bool    `json:"pending,omitempty" mapstructure:"pending"`
	PremiumSince               any      `json:"premium_since" mapstructure:"premium_since"`
}
