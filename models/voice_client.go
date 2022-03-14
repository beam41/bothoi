package models

import "sync"

type SongItem struct {
	YtID        string
	Title       string
	Duration    string
	RequesterID string
}

type SongItemWData struct {
	YtID         string
	Title        string
	Duration     string
	RequesterID  string
	SongData     []byte
	DownloadLock sync.Mutex
	Downloading  bool
}

type VoiceServer struct {
	Token    string `json:"token" mapstructure:"token"`
	GuildID  string `json:"guild_id" mapstructure:"guild_id"`
	Endpoint string `json:"endpoint" mapstructure:"endpoint"`
}
