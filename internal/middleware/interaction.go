package middleware

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/internal/utils"
)

// Type alias for the interaction event signature.
type interactionEvent func(session *discordgo.Session, interaction *discordgo.InteractionCreate)

// Prevent the bot from responding to itself.
func InteractionCreate(next interactionEvent) interactionEvent {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		user := utils.GetDiscordUser(interaction)
		if user.ID == session.State.User.ID {
			return
		}

		next(session, interaction)
	}
}
