package repo

import (
	"bothoi/models/db_models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strings"
	"testing"
	"time"
)

func Test_GetQueueSize(t *testing.T) {
	db_, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return
	}
	db = db_
	err = db.AutoMigrate(&db_models.Song{})
	if err != nil {
		return
	}
	db.Create(&db_models.Song{
		Id:      1,
		GuildId: 12,
		Title:   "Not This",
	})
	db.Delete(&db_models.Song{}, 1)
	db.Create(&db_models.Song{
		GuildId: 12,
		Title:   "This",
	})
	db.Create(&db_models.Song{
		GuildId: 12,
		Title:   "This",
	})
	db.Create(&db_models.Song{
		GuildId: 13,
		Title:   "Thats",
	})

	result := GetQueueSize(12)
	expected := int64(2)

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	} else {
		t.Logf("expected: %v, got: %v", expected, result)
	}
}

func Test_GetNextSong(t *testing.T) {
	db_, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return
	}
	db = db_
	err = db.AutoMigrate(&db_models.Song{})
	if err != nil {
		return
	}
	db.Create(&db_models.Song{
		Id:          1,
		GuildId:     12,
		Title:       "Not This",
		RequestedAt: time.Unix(1000, 0),
	})
	db.Delete(&db_models.Song{}, 1)
	db.Create(&db_models.Song{
		GuildId:     12,
		Title:       "That",
		RequestedAt: time.Unix(1000, 0),
	})
	db.Create(&db_models.Song{
		GuildId:     12,
		Title:       "This",
		RequestedAt: time.Unix(1, 0),
	})
	db.Create(&db_models.Song{
		GuildId:     13,
		Title:       "That",
		RequestedAt: time.Unix(1000, 0),
	})

	result, _, _ := GetNextSong(12)
	expected := "This"

	if result.Title != expected {
		t.Errorf("expected: %v, got: %v", expected, result.Title)
	} else {
		t.Logf("expected: %v, got: %v", expected, result.Title)
	}
}

func Test_GetSongQueue(t *testing.T) {
	db_, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		return
	}
	db = db_
	err = db.AutoMigrate(&db_models.Song{})
	if err != nil {
		return
	}
	db.Create(&db_models.Song{
		Id:          1,
		GuildId:     12,
		Title:       "Not This",
		RequestedAt: time.Unix(1000, 0),
	})
	db.Delete(&db_models.Song{}, 1)
	db.Create(&db_models.Song{
		GuildId:     12,
		Title:       "4",
		RequestedAt: time.Unix(4, 0),
	})
	db.Create(&db_models.Song{
		GuildId:     12,
		Title:       "1",
		RequestedAt: time.Unix(1, 0),
	})
	db.Create(&db_models.Song{
		GuildId:     12,
		Title:       "2",
		RequestedAt: time.Unix(2, 0),
	})
	db.Create(&db_models.Song{
		GuildId:     12,
		Title:       "3",
		RequestedAt: time.Unix(3, 0),
	})
	db.Create(&db_models.Song{
		GuildId:     12,
		Title:       "5",
		RequestedAt: time.Unix(5, 0),
	})
	db.Create(&db_models.Song{
		GuildId:     13,
		Title:       "That",
		RequestedAt: time.Unix(1000, 0),
	})

	resultSlice := GetSongQueue(12, 4)
	resultSliceStr := make([]string, len(resultSlice))
	for i, song := range resultSlice {
		resultSliceStr[i] = song.Title
	}

	result := strings.Join(resultSliceStr, ",")

	expected := "1,2,3,4"

	if result != expected {
		t.Errorf("expected: %v, got: %v", expected, result)
	} else {
		t.Logf("expected: %v, got: %v", expected, result)
	}
}
