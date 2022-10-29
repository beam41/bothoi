package db_models

import "bothoi/models/types"

type Channel struct {
	ID       types.Snowflake `gorm:"primaryKey,type:INTEGER"`
	Type     uint8
	Position uint16
	Name     *string
	Bitrate  uint32
	GuildID  types.Snowflake `gorm:"type:INTEGER"`
	OwnerID  types.Snowflake `gorm:"type:INTEGER"`

	VoiceState []VoiceState `gorm:"foreignKey:ChannelID"`
	Songs      []Song       `gorm:"foreignKey:RequestChannelID"`
}
