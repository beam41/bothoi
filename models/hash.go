package models

type Hash struct {
	Hash    string `json:"hash" mapstructure:"hash"`
	Omitted bool   `json:"omitted" mapstructure:"omitted"`
}
