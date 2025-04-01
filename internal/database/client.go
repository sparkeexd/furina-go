package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgreSQL client.
type DB struct {
	pg *pgxpool.Pool
}

// Returns a PostgreSQL client.
func DatabaseClient() (*DB, error) {
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	client := &DB{pg: db}
	return client, nil
}
