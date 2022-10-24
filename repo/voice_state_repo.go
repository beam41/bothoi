package repo

import (
	"bothoi/models/db_models"
	"bothoi/models/discord_models"
	"bothoi/models/types"
	"gorm.io/gorm/clause"
)

func UpsertVoiceState(voiceState *discord_models.VoiceState) {
	dbVoiceStates := db_models.VoiceState{
		UserId:                  voiceState.UserId,
		GuildId:                 voiceState.GuildId,
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
	db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&dbVoiceStates)
}

func GetChannelIdByUserIdAndGuildId(userId, guildId types.Snowflake) *types.Snowflake {
	var voiceState db_models.VoiceState
	db.Select("ChannelId").Find(&voiceState, map[string]interface{}{"user_id": userId, "guild_id": guildId})
	return voiceState.ChannelId
}
