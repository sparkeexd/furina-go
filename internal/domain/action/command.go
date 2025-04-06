package action

import (
	"github.com/bwmarrin/discordgo"
)

// Slash command structure holding the application command and its respective handlers to be added and accessed by the bot.
type Command struct {
	Command *discordgo.ApplicationCommand
	Handler CommandHandler
}

// Discord bot command interaction handler.
type CommandHandler func(session *discordgo.Session, interaction *discordgo.InteractionCreate)

// Service's slash commands to be registered.
type CommandService interface {
	Commands() map[string]Command
}

// Create a new slash command.
func NewCommand(command *discordgo.ApplicationCommand, handler CommandHandler) Command {
	return Command{command, handler}
}
