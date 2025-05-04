package entity

import (
	"time"
)

// Model for "users" table.
// GuildID is the Discord server ID.
type User struct {
	ID         int
	GuildID    int
	LtokenV2   string
	LtmidV2    string
	LtuidV2    string
	ModifiedAt time.Time
	CreatedAt  time.Time
}

// Model for "accounts" table.
type Account struct {
	ID         int
	Game       string
	DiscordID  int
	ModifiedAt time.Time
	CreatedAt  time.Time
}
