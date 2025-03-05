package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/sparkeexd/mimo/internal/models"
)

var (
	// Discord bot session.
	session *discordgo.Session

	// Discord bot parameters.
	botToken string
)

// Create discord bot session.
func CreateSession() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading configs: %v", err)
	}

	botToken = os.Getenv("BOT_TOKEN")
	session, err = discordgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
}

// Register the slash commands.
// Requires reloading Discord client to view the changes.
func AddCommands(commandGroups ...map[string]models.Command) {
	var commandsToRegister []*discordgo.ApplicationCommand

	for _, commands := range commandGroups {
		session.AddHandler(
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
	session.ApplicationCommandBulkOverwrite(session.State.User.ID, "", commandsToRegister)
}

// Close the session.
func CloseSession() {
	session.Close()
}
