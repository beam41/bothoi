package models

type Interaction struct {
	ApplicationID string          `json:"application_id" mapstructure:"application_id"`
	ChannelID     string          `json:"channel_id" mapstructure:"channel_id"`
	Data          InteractionData `json:"data" mapstructure:"data"`
	GuildID       string          `json:"guild_id" mapstructure:"guild_id"`
	GuildLocale   string          `json:"guild_locale" mapstructure:"guild_locale"`
	ID            string          `json:"id" mapstructure:"id"`
	Locale        string          `json:"locale" mapstructure:"locale"`
	Member        GuildMember     `json:"member" mapstructure:"member"`
	Token         string          `json:"token" mapstructure:"token"`
	Type          int64           `json:"type" mapstructure:"type"`
	Version       int64           `json:"version" mapstructure:"version"`
}

type InteractionData struct {
	ID      string              `json:"id" mapstructure:"id"`
	Name    string              `json:"name" mapstructure:"name"`
	Options []InteractionOption `json:"options" mapstructure:"options"`
	Type    int64               `json:"type" mapstructure:"type"`
}

type InteractionOption struct {
	Name  string      `json:"name" mapstructure:"name"`
	Type  int64       `json:"type" mapstructure:"type"`
	Value interface{} `json:"value" mapstructure:"value"`
}
