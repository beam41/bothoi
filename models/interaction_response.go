package models

type InteractionResponse struct {
	Type int8                    `json:"type"`
	Data InteractionResponseData `json:"data"`
}

type InteractionResponseData struct {
	Tts             *bool           `json:"tts,omitempty"`
	Content         *string         `json:"content,omitempty"`
	Embeds          []Embed         `json:"embeds,omitempty"`
	Flags           *int64          `json:"flags,omitempty"`
	AllowedMentions *AllowedMention `json:"allowed_mentions,omitempty"`
}
