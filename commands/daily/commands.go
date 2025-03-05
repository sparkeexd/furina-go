package daily

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/internal/models"
	"github.com/sparkeexd/mimo/internal/network"
)

var (
	// Command names.
	dailyCommandName = "daily"

	// Commands.
	Commands = map[string]models.Command{
		dailyCommandName: models.NewCommand(&dailyCommand, dailyCommandHandler),
	}
)

// Daily check-in claim command.
var dailyCommand = discordgo.ApplicationCommand{
	Name:        dailyCommandName,
	Description: "Command for Genshin daily check-in.",
}

// Perform Genshin Impact daily check-in on HoYoLab.
func dailyCommandHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	ltokenV2 := os.Getenv("LTOKEN_V2")
	ltmidV2 := os.Getenv("LTMID_V2")
	ltuidV2 := os.Getenv("LTUID_V2")
	cookie := network.NewCookie(ltokenV2, ltmidV2, ltuidV2)

	daily := NewDailyReward(Hk4eEndpoint, GenshinEventId, GenshinActId, GenshinSignGame)
	res, err := daily.Claim(cookie)

	message := fmt.Sprintf("You have successfully checked in, %s!", interaction.Member.Mention())

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
