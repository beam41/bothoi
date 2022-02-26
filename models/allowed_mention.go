package models

type AllowedMention struct {
	Parse []string `json:"parse,omitempty"`
	Users []string `json:"users,omitempty"`
}
