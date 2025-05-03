package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sparkeexd/mimo/internal/domain/entity"
)

// Repository for handling Discord users in the database.
type DiscordUserRepository struct {
	db *pgxpool.Pool
}

// Create a new HoyoLAB user repository.
func NewDiscordUserRepository(db *pgxpool.Pool) DiscordUserRepository {
	return DiscordUserRepository{db: db}
}

// Get discord user by ID.
func (repository DiscordUserRepository) GetByDiscordID(discordID int) (entity.DiscordUser, error) {
	query := `
		SELECT du.*
		FROM discord_users du
		WHERE du.id = @discordID
	`
	args := pgx.NamedArgs{
		"discordID": discordID,
	}

	var users entity.DiscordUser
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return users, err
	}

	users, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[entity.DiscordUser])

	return users, err
}

// Lists a batch of Discord users.
// Starts after the specificied Discord ID, limited by the given batch size.
func (repository DiscordUserRepository) ListDiscordUsers(offsetDiscordID int, limit int) ([]entity.DiscordUser, error) {
	query := `
		SELECT du.*
		FROM discord_users du
		WHERE du.id > @discordID
		ORDER BY du.id
		LIMIT @limit;
	`
	args := pgx.NamedArgs{
		"discordID": offsetDiscordID,
		"limit":     limit,
	}

	var users []entity.DiscordUser
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return users, err
	}

	users, err = pgx.CollectRows(rows, pgx.RowToStructByName[entity.DiscordUser])

	return users, err
}
