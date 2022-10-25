package db_models

import (
	"bothoi/models/types"
	"gorm.io/gorm"
	"time"
)

type Song struct {
	Id          uint32          `gorm:"primaryKey"`
	GuildId     types.Snowflake `gorm:"type:INTEGER"`
	RequesterId types.Snowflake `gorm:"type:INTEGER"`
	RequestedAt time.Time
	YtId        string
	Title       string
	Duration    uint32
	Seek        uint32
	Playing     bool
	Deleted     gorm.DeletedAt
}
