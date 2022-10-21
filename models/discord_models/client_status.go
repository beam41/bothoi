package discord_models

type ClientStatus struct {
	Desktop string `json:"desktop" mapstructure:"desktop"`
	Mobile  string `json:"mobile" mapstructure:"mobile"`
	Web     string `json:"web" mapstructure:"web"`
}
