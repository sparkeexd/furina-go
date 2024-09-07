package daily

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/hoyoapi/middleware"
	"github.com/sparkeexd/mimo/bot"
	"github.com/sparkeexd/mimo/commands"
)

var (
	// Command names.
	dailyCommandName = "daily"
	ltokenV2         = "ltokenv2"
	ltmidV2          = "ltmidv2"
	ltuidV2          = "ltuidv2"

	// Commands.
	Commands = map[string]commands.Command{
		dailyCommandName: {
			Command: &dailyCommand,
			Handler: dailyCommandHandler,
		},
	}
)

// Daily check-in claim command.
var dailyCommand = discordgo.ApplicationCommand{
	Name:        dailyCommandName,
	Description: "Command to claim HoyoLab daily check-in.",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Name:        ltokenV2,
			Description: "The ltokenV2 cookie",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
		{
			Name:        ltmidV2,
			Description: "The ltmidV2 cookie",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
		{
			Name:        ltuidV2,
			Description: "The ltuidV2 cookie",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		},
	},
}

// The bot will do daily check-in on HoYoLab.
// Calls the user by their display name or server nickname if present, otherwise defaults to their username.
func dailyCommandHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := make(map[string]string, len(interaction.ApplicationCommandData().Options))
	for _, opt := range interaction.ApplicationCommandData().Options {
		options[opt.Name] = opt.StringValue()
	}

	cookie := middleware.NewCookie(options[ltokenV2], options[ltmidV2], options[ltuidV2])

	client := bot.ZenlessClient
	client.Handler.Cookie = cookie

	message := "You have successfully checked in!"
	errorMessage := "Could not claim daily rewards"

	res, err := bot.ZenlessClient.Daily.Claim()

	if err != nil {
		message = fmt.Sprintf("%v: %v", errorMessage, err)
	} else if res.Retcode != 0 {
		message = fmt.Sprintf("%v: %v", errorMessage, res.Message)
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
