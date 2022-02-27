package models

type User struct {
	Avatar        *string `json:"avatar" mapstructure:"avatar"`
	Discriminator string  `json:"discriminator" mapstructure:"discriminator"`
	ID            string  `json:"id" mapstructure:"id"`
	Username      string  `json:"username" mapstructure:"username"`
	Bot           *bool   `json:"bot,omitempty" mapstructure:"bot"`
}
