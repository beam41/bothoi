package db_models

import (
	"bothoi/models/types"
	"time"
)

type Guild struct {
	Id          types.Snowflake `gorm:"primaryKey,type:INTEGER"`
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
	OwnerId     types.Snowflake `gorm:"type:INTEGER"`
	VoiceStates []VoiceState    `gorm:"constraint:OnDelete:CASCADE;foreignKey:GuildId"`
	Members     []GuildMember   `gorm:"constraint:OnDelete:CASCADE;foreignKey:GuildId"`
	Channels    []Channel       `gorm:"constraint:OnDelete:CASCADE;foreignKey:GuildId"`
	Songs       []Song          `gorm:"foreignKey:GuildId"`
}
