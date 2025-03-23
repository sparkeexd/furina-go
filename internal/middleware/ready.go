package middleware

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Log when the bot is ready to start receiving commands.
func Ready(session *discordgo.Session, ready *discordgo.Ready) {
	log.Printf("Logged in as: %v#%v", session.State.User.Username, session.State.User.Discriminator)
}
