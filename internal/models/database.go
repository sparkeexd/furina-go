package models

import (
	"os"
	"sync"

	"github.com/nedpals/supabase-go"
)

var (
	// Supabase client singleton.
	instance *database

	// Mutex to initialize singleton.
	mutex = &sync.Mutex{}
)

// Supabase client.
type database struct {
	client *supabase.Client
}

// Response structure from "users" table.
type UserResponse struct {
	DiscordID int    `json:"discord_id"`
	LtokenV2  string `json:"ltoken_v2"`
	LtmidV2   string `json:"ltmid_v2"`
	LtuidV2   string `json:"ltuid_v2"`
	CreatedAt string `json:"created_at"`
}

// Returns a Supabase client singleton.
func DatabaseClient() *database {
	mutex.Lock()
	defer mutex.Unlock()

	if instance == nil {
		projectURL := os.Getenv("SUPABASE_PROJECT_URL")
		apiKey := os.Getenv("SUPABASE_PUBLISHABLE_KEY")
		client := supabase.CreateClient(projectURL, apiKey)

		instance = &database{client: client}
	}

	return instance
}

// Get user's ltoken_v2, ltmid_v2, and ltuid_v2 tokens from the database.
func (database *database) GetUser(discordID string) (UserResponse, error) {
	var user UserResponse

	client := database.client
	err := client.DB.
		From("users").
		Select().
		Single().
		Eq("discord_id", discordID).
		Execute(&user)

	return user, err
}
