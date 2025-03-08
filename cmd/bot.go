package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/commands/daily"
	"github.com/sparkeexd/mimo/commands/hello"
	"github.com/sparkeexd/mimo/internal/models"
)

// Discord bot.
type Bot struct {
	Token    string
	Session  *discordgo.Session
	Commands []map[string]models.Command
	Status   string
}

// Create a new Discord bot.
func NewBot() *Bot {
	token := os.Getenv("BOT_TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	bot := &Bot{
		Token:   token,
		Session: session,
		Commands: []map[string]models.Command{
			hello.Commands,
			daily.Commands,
		},
		Status: "Active",
	}

	bot.Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	return bot
}

// Start Discord bot.
func (bot *Bot) Start() {
	// Start healthcheck server in a separate goroutine.
	server := NewServer(bot)
	go server.StartServer()

	log.Println("Creating discord bot session...")
	err := bot.Session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	bot.AddCommands(bot.Commands...)

	// Event listener to stop the bot.
	log.Println("Bot is now running! Press Ctrl+C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Closing discord bot session...")
	bot.Status = "Inactive"
	bot.Session.Close()
}

// Register the slash commands.
// Requires reloading Discord client to view the changes.
func (bot *Bot) AddCommands(commandGroups ...map[string]models.Command) {
	var commandsToRegister []*discordgo.ApplicationCommand

	for _, commands := range commandGroups {
		bot.Session.AddHandler(
			func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
				if command, exists := commands[interaction.ApplicationCommandData().Name]; exists {
					command.Handler(session, interaction)
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
