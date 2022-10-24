package db_models

import "bothoi/models/types"

type VoiceState struct {
	SessionId               string          `gorm:"primaryKey"`
	UserId                  types.Snowflake `gorm:"index,type:INTEGER"`
	GuildId                 types.Snowflake `gorm:"index,type:INTEGER"`
	Deaf                    bool
	Mute                    bool
	SelfDeaf                bool
	SelfMute                bool
	SelfStream              bool
	SelfVideo               bool
	Suppress                bool
	RequestToSpeakTimestamp *types.ISOTimeStamp
	ChannelId               *types.Snowflake `gorm:"type:text"`
}
