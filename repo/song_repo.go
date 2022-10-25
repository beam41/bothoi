package repo

import (
	"bothoi/models/db_models"
	"bothoi/models/types"
	"errors"
	"gorm.io/gorm"
	"time"
)

func AddSongToQueue(guildId, requesterId types.Snowflake, ytId, title string, duration uint32) error {
	result := db.Create(&db_models.Song{
		GuildId:     guildId,
		RequesterId: requesterId,
		RequestedAt: time.Now(),
		YtId:        ytId,
		Title:       title,
		Duration:    duration,
	})
	return result.Error
}

func GetQueueSize(guildId types.Snowflake) (size int64) {
	db.Model(&db_models.Song{}).Where(map[string]interface{}{"guild_id": guildId}).Count(&size)
	return
}

func GetSongQueue(guildId types.Snowflake, max int) (songs []db_models.Song) {
	db.
		Model(&db_models.Song{}).
		Select("RequesterId", "YtId", "Title", "Duration", "Playing").
		Where(map[string]interface{}{"guild_id": guildId}).
		Order("requested_at").
		Limit(max).
		Find(&songs)
	return
}

func DeleteSong(id uint32) error {
	result := db.Delete(&db_models.Song{}, id)
	return result.Error
}

func DeleteSongsInGuild(guildId types.Snowflake) error {
	result := db.
		Where(map[string]interface{}{"guild_id": guildId}).
		Delete(&db_models.Song{})
	return result.Error
}

func GetNextSong(guildId types.Snowflake) (db_models.Song, bool, error) {
	var song db_models.Song
	result := db.
		Where(map[string]interface{}{"guild_id": guildId}).
		Order("requested_at").
		Take(&song)
	return song, errors.Is(result.Error, gorm.ErrRecordNotFound), result.Error
}

func UpdateSongPlaying(id uint32) error {
	song := db_models.Song{
		Id: id,
	}
	result := db.Model(&song).Update("playing", true)
	return result.Error
}
