package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sparkeexd/mimo/internal/domain/entity"
)

// Repository for handling Discord users in the database.
type UserRepository struct {
	db *pgxpool.Pool
}

// Create a new user repository.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return UserRepository{db: db}
}

// Get user by ID.
func (repository UserRepository) GetByDiscordID(discordID int) (entity.User, error) {
	query := `
		SELECT u.*
		FROM users u
		WHERE u.id = @discordID
	`
	args := pgx.NamedArgs{
		"discordID": discordID,
	}

	var users entity.User
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return users, err
	}

	users, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[entity.User])

	return users, err
}

// Lists a batch of users.
// Starts after the specificied Discord ID, limited by the given batch size.
func (repository UserRepository) ListDiscordUsers(offsetDiscordID int, limit int) ([]entity.User, error) {
	query := `
		SELECT u.*
		FROM users u
		WHERE u.id > @discordID
		ORDER BY u.id
		LIMIT @limit;
	`
	args := pgx.NamedArgs{
		"discordID": offsetDiscordID,
		"limit":     limit,
	}

	var users []entity.User
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return users, err
	}

	users, err = pgx.CollectRows(rows, pgx.RowToStructByName[entity.User])

	return users, err
}
