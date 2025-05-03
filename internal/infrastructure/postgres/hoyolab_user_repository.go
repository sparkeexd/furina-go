package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sparkeexd/mimo/internal/domain/entity"
)

// Repository for handling HoyoLAB users in the database.
type HoyolabUserRepository struct {
	db *pgxpool.Pool
}

// Create a new HoyoLAB user repository.
func NewHoyolabUserRepository(db *pgxpool.Pool) HoyolabUserRepository {
	return HoyolabUserRepository{db: db}
}

// Get user's tokens by Discord ID.
func (repository HoyolabUserRepository) GetByDiscordID(discordID int) (entity.HoyolabUser, error) {
	query := `
		SELECT *
		FROM hoyolab_users hu
		WHERE hu.discord_id = @discordID;
	`
	args := pgx.NamedArgs{"discordID": discordID}

	var user entity.HoyolabUser
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return user, err
	}

	user, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[entity.HoyolabUser])
	return user, err
}
