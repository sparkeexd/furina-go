package models

import (
	"os"
	"sync"

	"github.com/sparkeexd/mimo/internal/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	// Database client singleton.
	instance *database

	// Mutex to initialize singleton.
	mutex = &sync.Mutex{}
)

// Database client.
type database struct {
	db *gorm.DB
}

// Model for "hoyolab_tokens" table.
type HoyolabToken struct {
	UserID    int `gorm:"primaryKey"`
	LtokenV2  string
	LtmidV2   string
	LtuidV2   string
	CreatedAt string
}

// Returns a PostgreSQL GORM client singleton.
func DatabaseClient() (*database, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if instance == nil {
		dsn := os.Getenv("DATABASE_URL")
		logLevel := utils.LogLevel()

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})

		if err != nil {
			return nil, err
		}

		instance = &database{db: db}
	}

	return instance, nil
}

// Get user's ltoken_v2, ltmid_v2, and ltuid_v2 tokens from the database.
func (database *database) HoyolabToken(userID int) (HoyolabToken, error) {
	var token HoyolabToken

	tx := database.db.First(&token, userID)
	return token, tx.Error
}
