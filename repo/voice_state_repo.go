package repo

import (
	"bothoi/models/db_models"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"gorm.io/gorm/clause"
)

func UpsertVoiceState(voiceState *discord_models.VoiceState) {
	dbVoiceStates := db_models.VoiceState{
		UserID:                  voiceState.UserID,
		GuildID:                 voiceState.GuildID,
		SessionID:               voiceState.SessionID,
		Deaf:                    voiceState.Deaf,
		Mute:                    voiceState.Mute,
		SelfDeaf:                voiceState.SelfDeaf,
		SelfMute:                voiceState.SelfMute,
		SelfStream:              voiceState.SelfStream,
		SelfVideo:               voiceState.SelfVideo,
		Suppress:                voiceState.Suppress,
		RequestToSpeakTimestamp: voiceState.RequestToSpeakTimestamp,
		ChannelID:               voiceState.ChannelID,
	}
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&dbVoiceStates)
}

func GetChannelIDByUserIDAndGuildID(userID, guildID types.Snowflake) *types.Snowflake {
	var voiceState db_models.VoiceState
	db.Select("ChannelID").Find(&voiceState, map[string]interface{}{"user_id": userID, "guild_id": guildID})
	return voiceState.ChannelID
}
