package entity

import "time"

// Model for "discord_users" table.
// GuildID is the Discord server ID.
type DiscordUser struct {
	ID        int
	GuildID   int
	CreatedAt time.Time
}

// Model for "hoyolab_users" table.
// ID is the HoYoLab account ID, which is the same value for the ltuid_v2 cookie.
type HoyolabUser struct {
	ID         int
	DiscordID  int
	LtokenV2   string
	LtmidV2    string
	ModifiedAt time.Time
	CreatedAt  time.Time
}

// Model for "game_users" table.
type GameUser struct {
	ID         int
	DiscordID  int
	GameID     int
	RegionID   int
	ModifiedAt time.Time
	CreatedAt  time.Time
}

// Model for "games" table.
type Game struct {
	ID        int
	Name      string
	CreatedAt time.Time
}

// Model for "game_regions" table.
type GameRegion struct {
	ID        int
	Name      string
	ResetTime time.Time
	CreatedAt time.Time
}
