package repo

import (
	"bothoi/models/db_models"
	"bothoi/models/types"
	"bothoi/util/player_util.go"
	"time"
)

func AddSongToQueue(guildID, requesterID, requestChannelID types.Snowflake, ytID, title string, duration, seek uint32, postMsgPlaying bool) error {
	result := db.Create(&db_models.Song{
		GuildID:          guildID,
		RequesterID:      requesterID,
		RequestChannelID: requestChannelID,
		RequestedAt:      time.Now(),
		YtID:             ytID,
		Title:            title,
		Duration:         duration,
		Seek:             seek,
		PostMsgPlaying:   postMsgPlaying,
	})
	return result.Error
}

func AddSongToQueueMultiple(guildID, requesterID, requestChannelID types.Snowflake, ytResult []player_util.Video) error {
	songs := make([]db_models.Song, len(ytResult))
	for i, result := range ytResult {
		songs[i] = db_models.Song{
			GuildID:          guildID,
			RequesterID:      requesterID,
			RequestChannelID: requestChannelID,
			RequestedAt:      time.Now(),
			YtID:             result.ID,
			Title:            result.Title,
			Duration:         uint32(result.Duration),
			PostMsgPlaying:   true,
		}
	}
	result := db.Create(&songs)
	return result.Error
}

func GetQueueSize(guildID types.Snowflake) (size int64) {
	db.Model(&db_models.Song{}).Where(map[string]interface{}{"guild_id": guildID}).Count(&size)
	return
}

func GetSongQueue(guildID types.Snowflake, offset, limit int) (songs []db_models.Song) {
	db.
		Model(&db_models.Song{}).
		Select("RequesterID", "YtID", "Title", "Duration", "Playing").
		Where(map[string]interface{}{"guild_id": guildID}).
		Order("requested_at").
		Offset(offset).
		Limit(limit).
		Find(&songs)
	return
}

func DeleteSong(id uint32) error {
	result := db.Delete(&db_models.Song{}, id)
	return result.Error
}

func DeleteSongsInGuild(guildID types.Snowflake) error {
	result := db.
		Where(map[string]interface{}{"guild_id": guildID}).
		Delete(&db_models.Song{})
	return result.Error
}

func GetNextSong(guildID types.Snowflake) (*db_models.Song, error) {
	var song db_models.Song
	result := db.
		Where(map[string]interface{}{"guild_id": guildID}).
		Order("requested_at").
		Limit(1).
		Find(&song)
	if result.RowsAffected == 0 {
		return nil, result.Error
	} else {
		return &song, result.Error
	}
}

func UpdateSongPlaying(id uint32) error {
	song := db_models.Song{
		ID: id,
	}
	result := db.Model(&song).Update("playing", true)
	return result.Error
}
