package daily

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/commands"
	"github.com/sparkeexd/mimo/internal/network"
)

var (
	// Command names.
	dailyCommandName = "daily"

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
	Description: "Command for Genshin daily check-in.",
}

// The bot will do daily check-in on HoYoLab.
// Calls the user by their display name or server nickname if present, otherwise defaults to their username.
func dailyCommandHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.
	options := make([]string, len(interaction.ApplicationCommandData().Options))
	for i, opt := range interaction.ApplicationCommandData().Options {
		options[i] = opt.Name
	}

	ltokenV2 := os.Getenv("LTOKEN_V2")
	ltmidV2 := os.Getenv("LTMID_V2")
	ltuidV2 := os.Getenv("LTUID_V2")
	cookie := network.NewCookie(ltokenV2, ltmidV2, ltuidV2)

	genshinDaily := NewDailyReward(Hk4eEndpoint, GenshinEventId, GenshinActId, GenshinSignGame)
	res, err := genshinDaily.Claim(cookie)

	message := "You have successfully checked in!"

	if err != nil {
		message = fmt.Sprint(err)
	} else if res.Retcode != 0 {
		message = res.Message
	}

	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
}
