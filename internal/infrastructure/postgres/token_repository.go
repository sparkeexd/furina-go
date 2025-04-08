package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Model for "tokens" table.
type Token struct {
	UserID    int       `db:"user_id"`
	LtokenV2  string    `db:"ltoken_v2"`
	LtmidV2   string    `db:"ltmid_v2"`
	LtuidV2   string    `db:"ltuid_v2"`
	CreatedAt time.Time `db:"created_at"`
}

// Repository for handling HoYoLab tokens in the database.
type TokenRepository struct {
	db *pgxpool.Pool
}

// Create a new token repository.
func NewTokenRepository(db *pgxpool.Pool) TokenRepository {
	return TokenRepository{db: db}
}

// Get user's ltoken_v2, ltmid_v2, and ltuid_v2 tokens from the database.
func (repository TokenRepository) GetByUserID(userID int) (Token, error) {
	query := `
		SELECT *
		FROM tokens t
		WHERE t.user_id = @userId;
	`
	args := pgx.NamedArgs{"userId": userID}

	var token Token
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return token, err
	}

	token, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[Token])
	return token, err
}

// Lists a batch of tokens starting from a specific userID, limited by the given batch size.
func (repository TokenRepository) ListByBatch(startUserID int, batchSize int) ([]Token, error) {
	query := `
		SELECT *
		FROM tokens t
		WHERE t.user_id > @userId
		ORDER BY t.user_id
		LIMIT @limit;
	`
	args := pgx.NamedArgs{
		"userId": startUserID,
		"limit":  batchSize,
	}

	var tokens []Token
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return tokens, err
	}

	tokens, err = pgx.CollectRows(rows, pgx.RowToStructByName[Token])

	return tokens, err
}
