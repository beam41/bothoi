package models

import (
	"bothoi/models/types"
	"sync"
)

type SongItem struct {
	YtId        string
	Title       string
	Duration    string
	RequesterId types.Snowflake
}

type SongItemWData struct {
	YtId         string
	Title        string
	Duration     string
	RequesterId  types.Snowflake
	SongData     []byte
	DownloadLock sync.Mutex
	Downloading  bool
}
