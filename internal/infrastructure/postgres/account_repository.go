package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sparkeexd/mimo/internal/domain/entity"
)

const (
	// Game enums.
	GenshinEnum  = "gi"
	StarRailEnum = "hsr"
	ZenlessEnum  = "zzz"

	// Game titles.
	GenshinImpact   = "Genshin Impact"
	HonkaiStarRail  = "Honkai: Star Rail"
	ZenlessZoneZero = "Zenless Zone Zero"
)

// Mapping of game enums to their game titles.
var gameTitles = map[string]string{
	GenshinEnum:  GenshinImpact,
	StarRailEnum: HonkaiStarRail,
	ZenlessEnum:  ZenlessZoneZero,
}

// Repository for handling HoYoverse accounts in the database.
type AccountRepository struct {
	db *pgxpool.Pool
}

// Create a new account repository.
func NewAccountRepository(db *pgxpool.Pool) AccountRepository {
	return AccountRepository{db: db}
}

// List HoYoverse game accounts based on Discord ID.
func (repository AccountRepository) ListByDiscordID(discordID int) ([]entity.Account, error) {
	query := `
		SELECT a.*
		FROM accounts a
		WHERE a.discord_id = @discordID
		ORDER BY a.game;
	`

	args := pgx.NamedArgs{
		"discordID": discordID,
	}

	var users []entity.Account
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return users, err
	}

	users, err = pgx.CollectRows(rows, pgx.RowToStructByName[entity.Account])
	return users, err
}

// Return the game title by the enum.
func (repository AccountRepository) GetGameTitle(gameID string) string {
	if name, exists := gameTitles[gameID]; exists {
		return name
	}

	return ""
}
