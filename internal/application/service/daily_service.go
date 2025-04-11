package service

import (
	"fmt"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"github.com/sparkeexd/mimo/internal/domain/action"
	"github.com/sparkeexd/mimo/internal/domain/network"
	"github.com/sparkeexd/mimo/internal/infrastructure/hoyolab"
	"github.com/sparkeexd/mimo/internal/infrastructure/postgres"
	"github.com/sparkeexd/mimo/pkg"
)

// Service that handles daily check-in commands.
type DailyService struct {
	DailyRepository       hoyolab.DailyRepository
	HoyolabUserRepository postgres.HoyolabUserRepository
}

// Create a new daily service.
func NewDailyService(dailyRepository hoyolab.DailyRepository, hoyolabUserRepository postgres.HoyolabUserRepository) DailyService {
	return DailyService{
		DailyRepository:       dailyRepository,
		HoyolabUserRepository: hoyolabUserRepository,
	}
}

// Service's slash commands to be registered.
func (service *DailyService) Commands() map[string]action.Command {
	return map[string]action.Command{
		"daily": action.NewCommand(
			&discordgo.ApplicationCommand{
				Name:        "daily",
				Description: "Command for Genshin daily check-in.",
			},
			service.DailyClaimCommandHandler,
		),
	}
}

// Service's cron jobs to be registered.
func (service *DailyService) Jobs(session *discordgo.Session) []action.CronJob {
	return []action.CronJob{
		action.NewCronJob(
			gocron.CronJob("0 16 * * *", false),
			gocron.NewTask(service.AutoDailyClaimTaskHandler, session),
		),
	}
}

// Perform Genshin Impact daily check-in on HoYoLAB.
func (service *DailyService) DailyClaimCommandHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	discordUser := pkg.GetDiscordUser(interaction)
	discordID, err := strconv.Atoi(discordUser.ID)
	if err != nil {
		content := "Invalid Discord user."
		pkg.InteractionResponseEditError(session, interaction.Interaction, err, content)
		return
	}

	user, err := service.HoyolabUserRepository.GetByDiscordID(discordID)
	if err != nil {
		content := "You are not registered yet, please register first."
		pkg.InteractionResponseEditError(session, interaction.Interaction, err, content)
		return
	}

	cookie := network.NewCookie(user.LtokenV2, user.LtmidV2, strconv.Itoa(user.ID))
	context := hoyolab.NewDailyRewardContext(hoyolab.Hk4eEndpoint, hoyolab.GenshinEventID, hoyolab.GenshinActID, hoyolab.GenshinSignGame)

	res, err := service.DailyRepository.Claim(cookie, context)
	if err != nil {
		content := "An internal error occurred while trying to check in. Please try again later."
		pkg.InteractionResponseEditError(session, interaction.Interaction, err, content)
		return
	}

	content := fmt.Sprintf("You have successfully checked in, %s!", discordUser.Mention())
	if res.Retcode != 0 {
		content = res.Message
	}

	session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}

// Task handler that automatically handles Genshin Impact daily check-in for all registered users.
func (service *DailyService) AutoDailyClaimTaskHandler(session *discordgo.Session) {
	startDiscordID := -1
	batchSize := 50

	for {
		users, err := service.HoyolabUserRepository.ListByBatch(startDiscordID, batchSize)
		if err != nil {
			log.Printf("Failed to list users: %v", err)
			return
		}

		// If no more users are found, exit the loop. This means all users have been processed.
		if len(users) == 0 {
			return
		}

		for _, user := range users {
			cookie := network.NewCookie(user.LtokenV2, user.LtmidV2, strconv.Itoa(user.ID))
			context := hoyolab.NewDailyRewardContext(hoyolab.Hk4eEndpoint, hoyolab.GenshinEventID, hoyolab.GenshinActID, hoyolab.GenshinSignGame)

			res, err := service.DailyRepository.Claim(cookie, context)

			content := "We have successfully checked in for you today!"
			if res.Retcode != 0 {
				content = res.Message
			} else if err != nil {
				// TODO: Add retry task handler.
				content = "An internal error occurred while trying to check in for you today."
			}

			channel, err := session.UserChannelCreate(strconv.Itoa(user.ID))
			if err != nil {
				log.Printf("Failed to send message to user channel %d: %v", user.ID, err)
				return
			}

			session.ChannelMessageSend(channel.ID, content)
		}

		// Start next batch from the last Discord ID in the current batch.
		startDiscordID = users[len(users)-1].ID
	}
}
