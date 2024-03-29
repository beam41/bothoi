package db_models

import (
	"bothoi/models/types"
	"time"
)

type GuildMember struct {
	UserID                     types.Snowflake `gorm:"primaryKey,type:INTEGER"`
	GuildID                    types.Snowflake `gorm:"primaryKey,type:INTEGER,"`
	Nick                       *string
	Avatar                     *string
	JoinedAt                   time.Time
	PremiumSince               *time.Time
	Deaf                       bool
	Mute                       bool
	Pending                    bool
	Permissions                string
	CommunicationDisabledUntil *time.Time
	// Roles                    []types.Snowflake `gorm:"type:text"`
}
