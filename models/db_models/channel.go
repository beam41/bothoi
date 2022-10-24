package db_models

import "bothoi/models/types"

type Channel struct {
	Id       types.Snowflake `gorm:"primaryKey,type:INTEGER"`
	Type     uint8
	Position uint16
	Name     *string
	Bitrate  uint32
	GuildId  types.Snowflake `gorm:"type:INTEGER"`
	OwnerId  types.Snowflake `gorm:"type:INTEGER"`

	VoiceState []VoiceState `gorm:"foreignKey:ChannelId"`
}
