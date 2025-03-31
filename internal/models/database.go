package models

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// Database singleton.
	instance *postgres

	// Mutex to initialize singleton.
	mutex = &sync.Mutex{}
)

// PostgreSQL client.
type postgres struct {
	db *pgxpool.Pool
}

// Model for "hoyolab_tokens" table.
type HoyolabToken struct {
	UserID    int       `db:"user_id"`
	LtokenV2  string    `db:"ltoken_v2"`
	LtmidV2   string    `db:"ltmid_v2"`
	LtuidV2   string    `db:"ltuid_v2"`
	CreatedAt time.Time `db:"created_at"`
}

// Returns a PostgreSQL GORM client singleton.
func DatabaseClient() (*postgres, error) {
	mutex.Lock()
	defer mutex.Unlock()

	if instance == nil {
		db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			return nil, err
		}

		instance = &postgres{db: db}
	}

	return instance, nil
}

// Get user's ltoken_v2, ltmid_v2, and ltuid_v2 tokens from the database.
func (pg *postgres) HoyolabToken(userID int) (HoyolabToken, error) {
	query := `
		SELECT *
		FROM hoyolab_tokens ht
		WHERE ht.user_id = @userId;
	`
	args := pgx.NamedArgs{"userId": userID}

	rows, err := pg.db.Query(context.Background(), query, args)
	if err != nil {
		var token HoyolabToken
		return token, err
	}

	token, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[HoyolabToken])
	return token, err
}
