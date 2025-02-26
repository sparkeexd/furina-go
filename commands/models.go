package commands

import "github.com/bwmarrin/discordgo"

// Slash command structure holding the application command and its respective handlers to be added and accessed by the bot.
type Command struct {
	Command *discordgo.ApplicationCommand
	Handler func(session *discordgo.Session, interaction *discordgo.InteractionCreate)
}
