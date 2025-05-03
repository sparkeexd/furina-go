package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sparkeexd/mimo/internal/domain/entity"
)

const (
	// Game title IDs.
	GenshinImpact   GameTitle = "gi"
	HonkaiStarRail  GameTitle = "hsr"
	ZenlessZoneZero GameTitle = "zzz"
)

// Game ID of `games` table.
type GameTitle string

// Repository for handling HoYoverse game users in the database.
type GameUserRepository struct {
	db *pgxpool.Pool
}

// Create a new game user repository.
func NewGameUserRepository(db *pgxpool.Pool) GameUserRepository {
	return GameUserRepository{db: db}
}

// List HoYoverse game accounts based on Discord user ID.
// A Discord user can only have 1 account per game.
func (repository GameUserRepository) ListByDiscordID(discordID int) ([]entity.GameUser, error) {
	query := `
		SELECT gu.*
		FROM game_users gu
		WHERE gu.discord_id = @discordID
		ORDER BY gu.game_id;
	`

	args := pgx.NamedArgs{
		"discordID": discordID,
	}

	var users []entity.GameUser
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return users, err
	}

	users, err = pgx.CollectRows(rows, pgx.RowToStructByName[entity.GameUser])
	return users, err
}

// Get HoYoverse game by ID.
func (repository GameUserRepository) GetGameByID(gameID int) (entity.Game, error) {
	query := `
		SELECT g.*
		FROM games g
		WHERE g.id = @gameID
	`

	args := pgx.NamedArgs{
		"gameID": gameID,
	}

	var game entity.Game
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return game, err
	}

	game, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Game])
	return game, err
}
