package repo

import (
	"bothoi/models/db_models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

func StartDb() {
	db_, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed to connect database")
	}
	db = db_

	// Migrate the schema
	err = db.AutoMigrate(
		&db_models.Guild{},
		&db_models.GuildMember{},
		&db_models.Channel{},
		&db_models.User{},
		&db_models.VoiceState{},
		&db_models.Song{},
	)
	if err != nil {
		log.Fatalln("Unable to migrate")
	}

	// clean up before get new value from GUILD_CREATE
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&db_models.Guild{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&db_models.GuildMember{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&db_models.Channel{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&db_models.VoiceState{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&db_models.Song{})
}
