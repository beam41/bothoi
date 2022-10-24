package repo

import (
	"bothoi/models/db_models"
	"bothoi/models/types"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func Test_GetChannelIdByUserIdAndGuildId(t *testing.T) {
	db_, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return
	}
	db = db_
	err = db.AutoMigrate(&db_models.VoiceState{})
	if err != nil {
		return
	}
	channelId1 := types.Snowflake(40)
	db.Create(&db_models.VoiceState{
		UserId:    12,
		GuildId:   23,
		ChannelId: &channelId1,
		SessionId: "a",
	})
	channelId2 := types.Snowflake(50)
	db.Create(&db_models.VoiceState{
		UserId:    12,
		GuildId:   23,
		ChannelId: &channelId2,
		SessionId: "b",
	})
	result := GetChannelIdByUserIdAndGuildId(12, 23)
	expected := types.Snowflake(40)

	if *result != expected {
		t.Errorf("expected: %v, got: %v", expected, *result)
	} else {
		t.Logf("expected: %v, got: %v", expected, *result)
	}
}
