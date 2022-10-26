package db_models

import "bothoi/models/types"

type VoiceState struct {
	SessionID               string          `gorm:"primaryKey"`
	UserID                  types.Snowflake `gorm:"index,type:INTEGER"`
	GuildID                 types.Snowflake `gorm:"index,type:INTEGER"`
	Deaf                    bool
	Mute                    bool
	SelfDeaf                bool
	SelfMute                bool
	SelfStream              bool
	SelfVideo               bool
	Suppress                bool
	RequestToSpeakTimestamp *types.ISOTimeStamp
	ChannelID               *types.Snowflake `gorm:"type:text"`
}
