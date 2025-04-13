package bot

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/internal/application/util"
	"github.com/sparkeexd/mimo/internal/domain/action"
)

// Log when the bot is ready to start receiving commands.
func (bot *Bot) Ready(session *discordgo.Session, ready *discordgo.Ready) {
	bot.Logger.Info(fmt.Sprintf("Logged in as: %v#%v", session.State.User.Username, session.State.User.Discriminator))
}

// Prevent the bot from responding to itself before handling the command.
func (bot *Bot) InteractionCreate(next action.CommandHandler) action.CommandHandler {
	return func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		user := util.GetDiscordUser(interaction)
		if user.ID == session.State.User.ID {
			return
		}

		next(session, interaction)
	}
}
