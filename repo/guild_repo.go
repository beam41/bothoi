package repo

import (
	"bothoi/models/db_models"
	"bothoi/models/discord_models"
	"gorm.io/gorm/clause"
	"time"
)

func UpsertGuild(guild *discord_models.GuildCreate) {
	joinedAt, _ := time.Parse(time.RFC3339, string(guild.JoinedAt))
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&db_models.Guild{
		Id:          guild.Id,
		Name:        guild.Name,
		Icon:        guild.Icon,
		IconHash:    guild.IconHash,
		Owner:       guild.Owner,
		MaxMembers:  guild.MaxMembers,
		Description: guild.Description,
		JoinedAt:    joinedAt,
		MemberCount: guild.MemberCount,
		Large:       guild.Large,
		Unavailable: guild.Unavailable,
		OwnerId:     guild.OwnerId,
	})

	users := make([]db_models.User, len(guild.Members))
	members := make([]db_models.GuildMember, len(guild.Members))
	for i, member := range guild.Members {
		// user
		users[i] = db_models.User{
			Id:            member.User.Id,
			Username:      member.User.Username,
			Discriminator: member.User.Discriminator,
			Avatar:        member.User.Avatar,
			Bot:           member.User.Bot,
			System:        member.User.System,
			MfaEnabled:    member.User.MfaEnabled,
			Banner:        member.User.Banner,
			AccentColor:   member.User.AccentColor,
			Locale:        member.User.Locale,
			Verified:      member.User.Verified,
			Email:         member.User.Email,
		}

		// member
		joinedAt, _ := time.Parse(time.RFC3339, string(member.JoinedAt))
		var premiumSince *time.Time
		if member.PremiumSince != nil {
			p, _ := time.Parse(time.RFC3339, string(*member.PremiumSince))
			premiumSince = &p
		}
		var communicationDisabledUntil *time.Time
		if member.PremiumSince != nil {
			c, _ := time.Parse(time.RFC3339, string(*member.CommunicationDisabledUntil))
			communicationDisabledUntil = &c
		}
		members[i] = db_models.GuildMember{
			UserId:  member.User.Id,
			GuildId: guild.Id,
			Nick:    member.Nick,
			Avatar:  member.Avatar,
			//Roles:                      member.Roles,
			JoinedAt:                   joinedAt,
			PremiumSince:               premiumSince,
			Deaf:                       member.Deaf,
			Mute:                       member.Mute,
			Pending:                    member.Pending,
			Permissions:                member.Permissions,
			CommunicationDisabledUntil: communicationDisabledUntil,
		}
	}
	if len(guild.Members) > 0 {
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&users)
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&members)
	}

	channels := make([]db_models.Channel, len(guild.Channels))
	for i, channel := range guild.Channels {
		channels[i] = db_models.Channel{
			Id:       channel.Id,
			Type:     channel.Type,
			Position: channel.Position,
			Name:     channel.Name,
			Bitrate:  channel.Bitrate,
			GuildId:  guild.Id,
			OwnerId:  channel.OwnerId,
		}
	}
	if len(guild.Channels) > 0 {
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&channels)
	}

	voiceStates := make([]db_models.VoiceState, len(guild.VoiceStates))
	for i, voiceState := range guild.VoiceStates {
		voiceStates[i] = db_models.VoiceState{
			UserId:                  voiceState.UserId,
			GuildId:                 guild.Id,
			SessionId:               voiceState.SessionId,
			Deaf:                    voiceState.Deaf,
			Mute:                    voiceState.Mute,
			SelfDeaf:                voiceState.SelfDeaf,
			SelfMute:                voiceState.SelfMute,
			SelfStream:              voiceState.SelfStream,
			SelfVideo:               voiceState.SelfVideo,
			Suppress:                voiceState.Suppress,
			RequestToSpeakTimestamp: voiceState.RequestToSpeakTimestamp,
			ChannelId:               voiceState.ChannelId,
		}
	}
	if len(guild.VoiceStates) > 0 {
		db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&voiceStates)
	}
}