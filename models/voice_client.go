package models

import (
	"bothoi/models/types"
)

type SongItem struct {
	YtId        string
	Title       string
	Duration    string
	RequesterId types.Snowflake
}
