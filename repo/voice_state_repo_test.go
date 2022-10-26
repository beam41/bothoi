package repo

import (
	"bothoi/models/db_models"
	"bothoi/models/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func Test_GetChannelIDByUserIDAndGuildID(t *testing.T) {
	db_, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return
	}
	db = db_
	err = db.AutoMigrate(&db_models.VoiceState{})
	if err != nil {
		return
	}
	channelID1 := types.Snowflake(40)
	db.Create(&db_models.VoiceState{
		UserID:    12,
		GuildID:   23,
		ChannelID: &channelID1,
		SessionID: "a",
	})
	channelID2 := types.Snowflake(50)
	db.Create(&db_models.VoiceState{
		UserID:    12,
		GuildID:   23,
		ChannelID: &channelID2,
		SessionID: "b",
	})
	result := GetChannelIDByUserIDAndGuildID(12, 23)
	expected := types.Snowflake(40)

	if *result != expected {
		t.Errorf("expected: %v, got: %v", expected, *result)
	} else {
		t.Logf("expected: %v, got: %v", expected, *result)
	}
}
