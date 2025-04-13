package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HoYoverse games.
type GameTitle string

// Game server regions.
type ServerRegion string

const (
	GenshinImpact   GameTitle = "Genshin Impact"
	HonkaiStarRail  GameTitle = "Honkai: Star Rail"
	ZenlessZoneZero GameTitle = "Zenless Zone Zero"

	Asia    ServerRegion = "Asia"
	America ServerRegion = "America"
	Europe  ServerRegion = "Europe"
)

// Model for "games" table.
type Game struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

// Model for "game_regions" table.
type GameRegion struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	ResetTime time.Time `db:"reset_time"`
	CreatedAt time.Time `db:"created_at"`
}

// Repository for handling HoYoverse game users in the database.
type GameRepository struct {
	db *pgxpool.Pool
}

// Create a new game repository.
func NewGameRepository(db *pgxpool.Pool) GameRepository {
	return GameRepository{
		db: db,
	}
}

// Returns all HoYoverse games.
func (repository GameRepository) GetGames() ([]Game, error) {
	query := `
		SELECT *
		FROM games g
		ORDER BY g.id;
	`

	var games []Game
	rows, err := repository.db.Query(context.Background(), query)
	if err != nil {
		return games, err
	}

	games, err = pgx.CollectRows(rows, pgx.RowToStructByName[Game])
	return games, err
}

// Returns all game server regions.
func (repository GameRepository) GetRegions() ([]GameRegion, error) {
	query := `
		SELECT *
		FROM game_regions gr
		ORDER BY gr.id;
	`

	var regions []GameRegion
	rows, err := repository.db.Query(context.Background(), query)
	if err != nil {
		return regions, err
	}

	regions, err = pgx.CollectRows(rows, pgx.RowToStructByName[GameRegion])
	return regions, err
}
