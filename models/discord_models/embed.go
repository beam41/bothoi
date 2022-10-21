package discord_models

import (
	"bothoi/models/types"
)

type Embed struct {
	Title       string              `json:"title,omitempty" mapstructure:"title"`
	Type        string              `json:"types,omitempty" mapstructure:"types"`
	Description string              `json:"description,omitempty" mapstructure:"description"`
	Url         string              `json:"url,omitempty" mapstructure:"url"`
	Timestamp   *types.ISOTimeStamp `json:"timestamp,omitempty" mapstructure:"timestamp"`
	Color       uint32              `json:"color,omitempty" mapstructure:"color"`
	Footer      *EmbedFooter        `json:"footer,omitempty" mapstructure:"footer"`
	Image       *EmbedImage         `json:"image,omitempty" mapstructure:"image"`
	Thumbnail   *EmbedThumbnail     `json:"thumbnail,omitempty" mapstructure:"thumbnail"`
	Video       *EmbedVideo         `json:"video,omitempty" mapstructure:"video"`
	Provider    *EmbedProvider      `json:"provider,omitempty" mapstructure:"provider"`
	Author      *EmbedAuthor        `json:"author,omitempty" mapstructure:"author"`
	Fields      []EmbedField        `json:"fields,omitempty" mapstructure:"fields"`
}

type EmbedFooter struct {
	Text         string `json:"text" mapstructure:"text"`
	IconUrl      string `json:"icon_url,omitempty" mapstructure:"icon_url"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty" mapstructure:"proxy_icon_url"`
}

type EmbedImage struct {
	Url      string `json:"url" mapstructure:"url"`
	ProxyUrl string `json:"proxy_url,omitempty" mapstructure:"proxy_url"`
	Height   uint32 `json:"height,omitempty" mapstructure:"height"`
	Width    uint32 `json:"width,omitempty" mapstructure:"width"`
}

type EmbedThumbnail struct {
	Url      string `json:"url" mapstructure:"url"`
	ProxyUrl string `json:"proxy_url,omitempty" mapstructure:"proxy_url"`
	Height   uint32 `json:"height,omitempty" mapstructure:"height"`
	Width    uint32 `json:"width,omitempty" mapstructure:"width"`
}

type EmbedVideo struct {
	Url      string `json:"url,omitempty" mapstructure:"url"`
	ProxyUrl string `json:"proxy_url,omitempty" mapstructure:"proxy_url"`
	Height   uint32 `json:"height,omitempty" mapstructure:"height"`
	Width    uint32 `json:"width,omitempty" mapstructure:"width"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty" mapstructure:"name"`
	Url  string `json:"url,omitempty" mapstructure:"url"`
}

type EmbedAuthor struct {
	Name         string `json:"name" mapstructure:"name"`
	Url          string `json:"url,omitempty" mapstructure:"url"`
	IconUrl      string `json:"icon_url,omitempty" mapstructure:"icon_url"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty" mapstructure:"proxy_icon_url"`
}

type EmbedField struct {
	Name   string `json:"name" mapstructure:"name"`
	Value  string `json:"value" mapstructure:"value"`
	Inline bool   `json:"inline,omitempty" mapstructure:"inline"`
}
