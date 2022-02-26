package models

import "time"

type Embed struct {
	Title       string          `json:"title,omitempty"`
	Type        string          `json:"type,omitempty"`
	Description string          `json:"description,omitempty"`
	Url         string          `json:"url,omitempty"`
	Timestamp   *time.Time      `json:"timestamp,omitempty"`
	Color       int32           `json:"color,omitempty"`
	Footer      *EmbedFooter    `json:"footer,omitempty"`
	Image       *EmbedImage     `json:"image,omitempty"`
	Thumbnail   *EmbedThumbnail `json:"thumbnail,omitempty"`
	Video       *EmbedVideo     `json:"video,omitempty"`
	Provider    *EmbedProvider  `json:"provider,omitempty"`
	Author      *EmbedAuthor    `json:"author,omitempty"`
	Fields      []EmbedField    `json:"fields,omitempty"`
}

type EmbedFooter struct {
	Text         string `json:"text"`
	IconUrl      string `json:"icon_url,omitempty"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type EmbedImage struct {
	Url      string `json:"url"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int32  `json:"height,omitempty"`
	Width    int32  `json:"width,omitempty"`
}

type EmbedThumbnail struct {
	Url       string `json:"url"`
	Proxy_url string `json:"proxy_url,omitempty"`
	Height    int32  `json:"height,omitempty"`
	Width     int32  `json:"width,omitempty"`
}

type EmbedVideo struct {
	Url      string `json:"url,omitempty"`
	ProxyUrl string `json:"proxy_url,omitempty"`
	Height   int32  `json:"height,omitempty"`
	Width    int32  `json:"width,omitempty"`
}

type EmbedProvider struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type EmbedAuthor struct {
	Name         string `json:"name"`
	Url          string `json:"url,omitempty"`
	IconUrl      string `json:"icon_url,omitempty"`
	ProxyIconUrl string `json:"proxy_icon_url,omitempty"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}
