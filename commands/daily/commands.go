package daily

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/mimo/internal/models"
	"github.com/sparkeexd/mimo/internal/network"
	"github.com/sparkeexd/mimo/internal/utils"
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
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	discordUser := utils.GetDiscordUser(interaction)
	userID, err := strconv.Atoi(discordUser.ID)
	if err != nil {
		content := "Invalid Discord user."
		utils.InteractionResponseEditError(session, interaction.Interaction, err, content)
		return
	}

	client, err := models.DatabaseClient()
	if err != nil {
		content := "Unable to connect to database."
		utils.InteractionResponseEditError(session, interaction.Interaction, err, content)
		return
	}

	token, err := client.HoyolabToken(userID)
	if err != nil {
		content := "You are not registered yet, please register first."
		utils.InteractionResponseEditError(session, interaction.Interaction, err, content)
		return
	}

	cookie := network.NewCookie(token.LtokenV2, token.LtmidV2, token.LtuidV2)
	daily := NewDailyReward(Hk4eEndpoint, GenshinEventID, GenshinActID, GenshinSignGame)
	res, err := daily.Claim(cookie)

	message := fmt.Sprintf("You have successfully checked in, %s!", discordUser.Mention())
	if err != nil {
		message = fmt.Sprint(err)
	} else if res.Retcode != 0 {
		message = res.Message
	}

	session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content: &message,
	})
}
