package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Model for "hoyolab_users" table.
// ID is the HoYoLab account ID, which is the same value for the ltuid_v2 cookie.
type HoyolabUser struct {
	ID         int       `db:"id"`
	DiscordID  int       `db:"discord_id"`
	LtokenV2   string    `db:"ltoken_v2"`
	LtmidV2    string    `db:"ltmid_v2"`
	ModifiedAt time.Time `db:"modified_at"`
	CreatedAt  time.Time `db:"created_at"`
}

// Repository for handling HoyoLAB users in the database.
type HoyolabUserRepository struct {
	db *pgxpool.Pool
}

// Create a new HoyoLAB user repository.
func NewHoyolabUserRepository(db *pgxpool.Pool) HoyolabUserRepository {
	return HoyolabUserRepository{
		db: db,
	}
}

// Get user's tokens by Discord ID.
func (repository HoyolabUserRepository) GetByDiscordID(discordID int) (HoyolabUser, error) {
	query := `
		SELECT *
		FROM hoyolab_users hu
		WHERE hu.discord_id = @discordID;
	`
	args := pgx.NamedArgs{"discordID": discordID}

	var user HoyolabUser
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return user, err
	}

	user, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[HoyolabUser])
	return user, err
}

// Lists a batch of users' tokens by their game account's region.
// Starts after the specificied Discord ID, limited by the given batch size.
func (repository HoyolabUserRepository) ListByRegionID(regionID int, offsetDiscordID int, batchSize int) ([]HoyolabUser, error) {
	query := `
		SELECT hu.*
		FROM hoyolab_users hu
		JOIN game_users gu ON hu.discord_id = gu.discord_id
		WHERE
			gu.region_id = @regionID
			AND hu.id > @discordID
		ORDER BY hu.id
		LIMIT @limit;
	`
	args := pgx.NamedArgs{
		"regionID":  regionID,
		"discordID": offsetDiscordID,
		"limit":     batchSize,
	}

	var users []HoyolabUser
	rows, err := repository.db.Query(context.Background(), query, args)
	if err != nil {
		return users, err
	}

	users, err = pgx.CollectRows(rows, pgx.RowToStructByName[HoyolabUser])

	return users, err
}
