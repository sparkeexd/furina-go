package middleware

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/internal/database"
	"github.com/sparkeexd/mimo/internal/models"
	"github.com/sparkeexd/mimo/internal/utils"
)

// Prevent the bot from responding to itself.
func InteractionCreate(next models.CommandHandler) models.CommandHandler {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *database.DB) {
		user := utils.GetDiscordUser(interaction)
		if user.ID == session.State.User.ID {
			return
		}

		next(session, interaction, db)
	}
}
