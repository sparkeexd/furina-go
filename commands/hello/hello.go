package hello

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/commands"
)

var (
	// Command names.
	helloCommandName = "hello"

	// Commands.
	Commands = map[string]commands.Command{
		helloCommandName: {
			Command: &helloCommand,
			Handler: helloCommandHandler,
		},
	}
)

// Hello command.
var helloCommand = discordgo.ApplicationCommand{
	Name:        helloCommandName,
	Description: "Basic hello greeting.",
}

// The bot will reply with a simple hello greeting to the user.
// Calls the user by their display name or server nickname if present, otherwise defaults to their username.
func helloCommandHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	member := interaction.Member

	name := member.User.Username
	globalName := member.User.GlobalName
	nickName := member.Nick

	if globalName != "" {
		name = globalName
	}

	if nickName != "" {
		name = nickName
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Hello there, %v!", name),
		},
	})
}
