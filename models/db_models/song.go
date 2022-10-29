package db_models

import (
	"bothoi/models/types"
	"gorm.io/gorm"
	"time"
)

type Song struct {
	ID               uint32          `gorm:"primaryKey"`
	GuildID          types.Snowflake `gorm:"type:INTEGER"`
	RequesterID      types.Snowflake `gorm:"type:INTEGER"`
	RequestChannelID types.Snowflake `gorm:"type:INTEGER"`
	RequestedAt      time.Time
	YtID             string
	Title            string
	Duration         uint32
	Seek             uint32
	Playing          bool
	PostMsgPlaying   bool
	Deleted          gorm.DeletedAt
}
