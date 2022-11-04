package db_models

import "bothoi/models/types"

type VoiceState struct {
	UserID                  types.Snowflake  `gorm:"primaryKey;not null;type:INTEGER"`
	GuildID                 types.Snowflake  `gorm:"primaryKey;not null;type:INTEGER"`
	ChannelID               *types.Snowflake `gorm:"type:text"`
	SessionID               string
	Deaf                    bool
	Mute                    bool
	SelfDeaf                bool
	SelfMute                bool
	SelfStream              bool
	SelfVideo               bool
	Suppress                bool
	RequestToSpeakTimestamp *types.ISOTimeStamp
}
