package models

type AllowedMention struct {
	Parse []string `json:"parse,omitempty" mapstructure:"parse"`
	Users []string `json:"users,omitempty" mapstructure:"users"`
}
