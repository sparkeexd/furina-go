package service

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/internal/application/util"
	"github.com/sparkeexd/mimo/internal/domain/action"
)

// Service that handles a basic ping command to the bot.
type PingService struct{}

// Create a new ping service.
func NewPingService() PingService {
	return PingService{}
}

// Service's slash commands to be registered.
func (service *PingService) Commands() map[string]action.Command {
	return map[string]action.Command{
		"ping": action.NewCommand(
			&discordgo.ApplicationCommand{
				Name:        "ping",
				Description: "Basic hello greeting.",
			},
			service.PingCommandHandler,
		),
	}
}

// Reply with a simple hello greeting to the user.
func (service *PingService) PingCommandHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	user := util.GetDiscordUser(interaction)
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Hello there, %v!", user.Mention()),
		},
	})
}
