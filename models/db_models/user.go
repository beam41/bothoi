package db_models

import "bothoi/models/types"

type User struct {
	ID            types.Snowflake `gorm:"primaryKey,type:INTEGER"`
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

	Guild       []Guild       `gorm:"foreignKey:OwnerID"`
	GuildMember []GuildMember `gorm:"foreignKey:UserID"`
	Channel     []Channel     `gorm:"foreignKey:OwnerID"`
	VoiceStates []VoiceState  `gorm:"foreignKey:UserID"`
	Song        []Song        `gorm:"foreignKey:RequesterID"`
}
