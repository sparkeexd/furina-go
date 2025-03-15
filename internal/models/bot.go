package models

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/internal/middleware"
)

// Discord bot.
type Bot struct {
	Token    string
	Session  *discordgo.Session
	Commands []map[string]Command
	Status   string
}

// Create a new Discord bot.
func NewBot(commands ...map[string]Command) *Bot {
	token := os.Getenv("BOT_TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	bot := &Bot{
		Token:    token,
		Session:  session,
		Commands: commands,
		Status:   "ACTIVE",
	}

	bot.Session.AddHandler(middleware.Ready)

	return bot
}

// Start Discord bot.
func (bot *Bot) Start() {
	log.Println("Creating discord bot session...")
	err := bot.Session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	bot.AddCommands()

	// Event listener to stop the bot.
	log.Println("Bot is now running! Press Ctrl+C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Closing discord bot session...")
	bot.Status = "INACTIVE"
	bot.Session.Close()
}

// Register the slash commands.
// Middleware is attached to each command to block interactions outside of the guild.
// Requires reloading Discord client to view the changes.
func (bot *Bot) AddCommands() {
	var commandsToRegister []*discordgo.ApplicationCommand

	for _, commands := range bot.Commands {
		bot.Session.AddHandler(
			func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
				if command, exists := commands[interaction.ApplicationCommandData().Name]; exists {
					middleware.InteractionCreate(command.Handler)(session, interaction)
				}
			},
		)

		for _, v := range commands {
			commandsToRegister = append(commandsToRegister, v.Command)
		}
	}

	// Overwrite all existing commands, which allow clearing out old commands.
	bot.Session.ApplicationCommandBulkOverwrite(bot.Session.State.User.ID, "", commandsToRegister)
}
