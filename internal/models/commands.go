package models

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/internal/database"
)

// Slash command structure holding the application command and its respective handlers to be added and accessed by the bot.
type Command struct {
	Command *discordgo.ApplicationCommand
	Handler CommandHandler
}

// Discord bot command interaction handler.
type CommandHandler func(session *discordgo.Session, interaction *discordgo.InteractionCreate, db *database.DB)

// Create a new slash command.
func NewCommand(command *discordgo.ApplicationCommand, handler CommandHandler) Command {
	return Command{command, handler}
}
