package storage

import (
	"log"
	"os"

	"github.com/nicolito128/nintendo-salta/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SqliteStorage representa una nueva conexión con la base de datos SQLite.
type PostgreSQLStorage struct {
	db *gorm.DB
}

// checking interface implementation
var _ Storage = &PostgreSQLStorage{}

func NewPostgreSQLStorage() *PostgreSQLStorage {
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_CONNECTION")), &gorm.Config{})
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

	return &PostgreSQLStorage{db}
}

func (ss *PostgreSQLStorage) DB() *gorm.DB {
	return ss.db
}
