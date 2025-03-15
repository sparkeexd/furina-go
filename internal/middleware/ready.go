package middleware

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Log when the bot is ready to start receiving commands.
func Ready(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
}
