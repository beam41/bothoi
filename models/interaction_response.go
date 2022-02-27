package models

type InteractionResponse struct {
	Type int8                    `json:"type" mapstructure:"type"`
	Data InteractionResponseData `json:"data" mapstructure:"data"`
}

type InteractionResponseData struct {
	Tts             *bool           `json:"tts,omitempty" mapstructure:"tts"`
	Content         *string         `json:"content,omitempty" mapstructure:"content"`
	Embeds          []Embed         `json:"embeds,omitempty" mapstructure:"embeds"`
	Flags           *int64          `json:"flags,omitempty" mapstructure:"flags"`
	AllowedMentions *AllowedMention `json:"allowed_mentions,omitempty" mapstructure:"allowed_mentions"`
}
