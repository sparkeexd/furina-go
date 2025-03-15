package middleware

import (
	"github.com/bwmarrin/discordgo"
)

// Type alias for the interaction event signature.
type interactionEvent func(session *discordgo.Session, interaction *discordgo.InteractionCreate)

// Prevent the bot from responding to itself and from DMs.
func InteractionCreate(next interactionEvent) interactionEvent {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// Check if the interaction was sent from a guild (server).
		// If GuildID is empty, the interaction was sent from a DM or group DM.
		if interaction.GuildID == "" {
			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This bot only accepts slash commands from a server, not from DMs.",
				},
			})
			return
		}

		// Ignore interactions from the bot itself.
		if interaction.Member.User.ID == session.State.User.ID {
			return
		}

		next(session, interaction)
	}
}
