package db

import (
	"github.com/charmbracelet/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"video_stream/pkg/model"
)

func Connect() *gorm.DB {
	log.Info("Connecting to database")

	db, err := gorm.Open(sqlite.Open("videos.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Connected to database")

	db.AutoMigrate(&model.Video{})

	return db
}
