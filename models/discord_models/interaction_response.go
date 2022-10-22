package discord_models

type InteractionResponse struct {
	Type uint8                   `json:"type" mapstructure:"type"`
	Data InteractionCallbackData `json:"data" mapstructure:"data"`
}

type InteractionCallbackData struct {
	Tts             *bool           `json:"tts,omitempty" mapstructure:"tts"`
	Content         *string         `json:"content,omitempty" mapstructure:"content"`
	Embeds          []Embed         `json:"embeds,omitempty" mapstructure:"embeds"`
	Flags           *uint8          `json:"flags,omitempty" mapstructure:"flags"`
	AllowedMentions *AllowedMention `json:"allowed_mentions,omitempty" mapstructure:"allowed_mentions"`
}
