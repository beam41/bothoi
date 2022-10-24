package discord_models

import "bothoi/models/types"

type User struct {
	Id            types.Snowflake `json:"id,string" mapstructure:"id"`
	Username      string          `json:"username" mapstructure:"username"`
	Discriminator string          `json:"discriminator" mapstructure:"discriminator"`
	Avatar        *string         `json:"avatar" mapstructure:"avatar"`
	Bot           bool            `json:"bot" mapstructure:"bot"`
	System        bool            `json:"system" mapstructure:"system"`
	MfaEnabled    bool            `json:"mfa_enabled" mapstructure:"mfa_enabled"`
	Banner        *string         `json:"banner" mapstructure:"banner"`
	AccentColor   *uint32         `json:"accent_color" mapstructure:"accent_color"`
	Locale        string          `json:"locale" mapstructure:"locale"`
	Verified      bool            `json:"verified" mapstructure:"verified"`
	Email         *string         `json:"email" mapstructure:"email"`
	Flags         uint32          `json:"flags" mapstructure:"flags"`
	PremiumType   uint8           `json:"premium_type" mapstructure:"premium_type"`
	PublicFlags   uint32          `json:"public_flags" mapstructure:"public_flags"`
}
