package hello

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/internal/models"
)

var (
	// Command names.
	helloCommandName = "hello"

	// Commands.
	Commands = map[string]models.Command{
		helloCommandName: models.NewCommand(&helloCommand, helloCommandHandler),
	}
)

// Hello command.
var helloCommand = discordgo.ApplicationCommand{
	Name:        helloCommandName,
	Description: "Basic hello greeting.",
}

// Reply with a simple hello greeting to the user.
// Calls the user by their display name or server nickname if present, otherwise defaults to their username.
func helloCommandHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Hello there, %v!", interaction.Member.Mention()),
		},
	})
}
