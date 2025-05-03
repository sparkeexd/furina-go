package service

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/internal/domain/entity"
	"github.com/sparkeexd/mimo/internal/infrastructure/discord"
)

// Service that handles a basic ping command to the bot.
type PingService struct {
	interactionRepository discord.InteractionRepository
}

// Create a new ping service.
func NewPingService(interactionRepository discord.InteractionRepository) PingService {
	return PingService{interactionRepository: interactionRepository}
}

// Service's slash commands to be registered.
func (service *PingService) Commands() map[string]entity.Command {
	return map[string]entity.Command{
		"hello": entity.NewCommand(
			&discordgo.ApplicationCommand{
				Name:        "hello",
				Description: "Say hello to Surintendante Chevalmarin.",
			},
			service.PingCommandHandler,
		),
	}
}

// Reply with a simple hello greeting to the user.
func (service *PingService) PingCommandHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	user := service.interactionRepository.GetDiscordUser(interaction)
	embed := service.interactionRepository.CreateEmbed().
		SetAuthor(user.Username, user.AvatarURL("")).
		SetTitle("Greetings").
		SetDescription(fmt.Sprintf("Hello there, %v!", user.Mention())).
		SetImage("https://media.tenor.com/NpR_anvK1woAAAAd/genshinimpact-furina.gif").
		SetThumbnail("https://media.tenor.com/EhXA2CCJ-QUAAAAj/furina.gif")

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}
