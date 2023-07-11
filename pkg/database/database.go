package database

import (
	"log"

	"github.com/nicolito128/nintendo-salta/pkg/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	DB = Open()
	log.Println("Database connected!")

	DB.AutoMigrate(models.User{})
	log.Println("User model migrated!")

	DB.AutoMigrate(models.Session{})
	log.Println("Session model migrated!")
}

func Open() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	return db
}
