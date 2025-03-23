package middleware

import (
	"github.com/bwmarrin/discordgo"
)

// Type alias for the interaction event signature.
type interactionEvent func(session *discordgo.Session, interaction *discordgo.InteractionCreate)

// Prevent the bot from responding to itself.
func InteractionCreate(next interactionEvent) interactionEvent {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if interaction.Member.User.ID == session.State.User.ID {
			return
		}

		next(session, interaction)
	}
}
