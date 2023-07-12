package storage

import (
	"log"

	"github.com/nicolito128/nintendo-salta/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SqliteStorage representa una nueva conexión con la base de datos SQLite.
type SqliteStorage struct {
	db *gorm.DB
}

// checking interface implementation
var _ Storage = &SqliteStorage{}

func NewSqliteStorage(storageName string) *SqliteStorage {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.Session{})
	if err != nil {
		log.Fatal(err)
	}

	return &SqliteStorage{db}
}

func (ss *SqliteStorage) DB() *gorm.DB {
	return ss.db
}
