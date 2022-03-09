package models

import "sync"

type SongItem struct {
	YtID         string
	Title        string
	Duration     int64
	RequesterID  string
}

type SongItemWData struct {
	YtID         string
	Title        string
	Duration     int64
	RequesterID  string
	SongData     []byte
	DownloadLock sync.Mutex
}

type VoiceServer struct {
	Token    string `json:"token" mapstructure:"token"`
	GuildID  string `json:"guild_id" mapstructure:"guild_id"`
	Endpoint string `json:"endpoint" mapstructure:"endpoint"`
}
