package db_models

import (
	"bothoi/models/types"
	"time"
)

type Guild struct {
	ID          types.Snowflake `gorm:"primaryKey,type:INTEGER"`
	Name        string
	Icon        *string
	IconHash    *string
	Owner       bool
	MaxMembers  uint32
	Description *string
	JoinedAt    time.Time
	MemberCount uint32
	Large       bool
	Unavailable bool
	OwnerID     types.Snowflake `gorm:"type:INTEGER"`
	VoiceStates []VoiceState    `gorm:"constraint:OnDelete:CASCADE;foreignKey:GuildID"`
	Members     []GuildMember   `gorm:"constraint:OnDelete:CASCADE;foreignKey:GuildID"`
	Channels    []Channel       `gorm:"constraint:OnDelete:CASCADE;foreignKey:GuildID"`
	Songs       []Song          `gorm:"foreignKey:GuildID"`
}
