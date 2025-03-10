package daily

import (
	"fmt"

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
	client := models.DatabaseClient()
	user, err := client.GetUser(interaction.Member.User.ID)

	if err != nil {
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "You are not registered yet, please register first.",
			},
		})
		return
	}

	cookie := network.NewCookie(user.LtokenV2, user.LtmidV2, user.LtuidV2)

	daily := NewDailyReward(Hk4eEndpoint, GenshinEventID, GenshinActID, GenshinSignGame)
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
