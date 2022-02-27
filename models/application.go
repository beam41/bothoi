package models

type Application struct {
	ID    string `json:"id" mapstructure:"id"`
	Flags int64  `json:"flags" mapstructure:"flags"`
}
