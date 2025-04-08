package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Model for "users" table.
type User struct {
	ID        int       `db:"id"`
	LtokenV2  string    `db:"ltoken_v2"`
	LtmidV2   string    `db:"ltmid_v2"`
	LtuidV2   string    `db:"ltuid_v2"`
	CreatedAt time.Time `db:"created_at"`
}

// Repository for handling Discord users in the database.
type UserRepository struct {
	db *pgxpool.Pool
}

// Create a new user repository.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return UserRepository{db: db}
}

// Get user's ltoken_v2, ltmid_v2, and ltuid_v2 tokens from the database.
func (repository UserRepository) GetByUserID(userID int) (User, error) {
	query := `
		SELECT *
		FROM users u
		WHERE u.id = @userId;
	`
	args := pgx.NamedArgs{"userId": userID}

	var user User
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return user, err
	}

	user, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	return user, err
}

// Lists a batch of users' tokens starting from a specific userID, limited by the given batch size.
func (repository UserRepository) ListByBatch(startUserID int, batchSize int) ([]User, error) {
	query := `
		SELECT *
		FROM users u
		WHERE u.id > @userId
		ORDER BY u.id
		LIMIT @limit;
	`
	args := pgx.NamedArgs{
		"userId": startUserID,
		"limit":  batchSize,
	}

	var users []User
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return users, err
	}

	users, err = pgx.CollectRows(rows, pgx.RowToStructByName[User])

	return users, err
}
