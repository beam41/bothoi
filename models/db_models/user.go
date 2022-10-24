package db_models

import "bothoi/models/types"

type User struct {
	Id            types.Snowflake `gorm:"primaryKey,type:INTEGER"`
	Username      string
	Discriminator string
	Avatar        *string
	Bot           bool
	System        bool
	MfaEnabled    bool
	Banner        *string
	AccentColor   *uint32
	Locale        string
	Verified      bool
	Email         *string
	Flags         uint32
	PremiumType   uint8
	PublicFlags   uint32

	Guild       []Guild       `gorm:"foreignKey:OwnerId"`
	GuildMember []GuildMember `gorm:"foreignKey:UserId"`
	Channel     []Channel     `gorm:"foreignKey:OwnerId"`
	VoiceStates []VoiceState  `gorm:"foreignKey:UserId"`
	Song        []Song        `gorm:"foreignKey:RequesterId"`
}
