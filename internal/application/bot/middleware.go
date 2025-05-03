package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/internal/domain/entity"
)

// Log when the bot is ready to start receiving commands.
func (bot *Bot) logReady(session *discordgo.Session, ready *discordgo.Ready) {
	bot.logger.Info(fmt.Sprintf("Logged in as: %v#%v", session.State.User.Username, session.State.User.Discriminator))
}

// Prevent the bot from responding to itself before handling the command.
func (bot *Bot) filterInteraction(next entity.CommandHandler, user *discordgo.User) entity.CommandHandler {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if user.ID == session.State.User.ID {
			return
		}

		next(session, interaction)
	}
}
