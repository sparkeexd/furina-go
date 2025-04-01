package database

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

// Model for "hoyolab_tokens" table.
type HoyolabToken struct {
	UserID    int       `db:"user_id"`
	LtokenV2  string    `db:"ltoken_v2"`
	LtmidV2   string    `db:"ltmid_v2"`
	LtuidV2   string    `db:"ltuid_v2"`
	CreatedAt time.Time `db:"created_at"`
}

// Get user's ltoken_v2, ltmid_v2, and ltuid_v2 tokens from the database.
func (db DB) GetHoyolabToken(userID int) (HoyolabToken, error) {
	query := `
		SELECT *
		FROM hoyolab_tokens ht
		WHERE ht.user_id = @userId;
	`
	args := pgx.NamedArgs{"userId": userID}

	var token HoyolabToken
	rows, err := db.pg.Query(context.Background(), query, args)
	if err != nil {
		return token, err
	}

	token, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[HoyolabToken])
	return token, err
}

// List all user tokens from the database.
func (db DB) ListHoyolabTokens() ([]HoyolabToken, error) {
	query := `
		SELECT *
		FROM hoyolab_tokens ht;
	`

	var tokens []HoyolabToken
	rows, err := db.pg.Query(context.Background(), query)
	if err != nil {
		return tokens, err
	}

	tokens, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[[]HoyolabToken])
	return tokens, err
}
