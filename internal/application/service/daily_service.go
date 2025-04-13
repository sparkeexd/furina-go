package service

import (
	"fmt"
	"log/slog"

	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"github.com/sparkeexd/mimo/internal/application/util"
	"github.com/sparkeexd/mimo/internal/domain/action"
	"github.com/sparkeexd/mimo/internal/domain/logger"
	"github.com/sparkeexd/mimo/internal/domain/network"
	"github.com/sparkeexd/mimo/internal/infrastructure/hoyolab"
	"github.com/sparkeexd/mimo/internal/infrastructure/postgres"
)

// Service that handles daily check-in commands.
type DailyService struct {
	DailyRepository       hoyolab.DailyRepository
	GameRepository        postgres.GameRepository
	HoyolabUserRepository postgres.HoyolabUserRepository
	logger                *logger.Logger
}

// Create a new daily service.
func NewDailyService(
	dailyRepository hoyolab.DailyRepository,
	hoyolabUserRepository postgres.HoyolabUserRepository,
	gameRepository postgres.GameRepository,
	logger *logger.Logger,
) DailyService {
	return DailyService{
		DailyRepository:       dailyRepository,
		HoyolabUserRepository: hoyolabUserRepository,
		GameRepository:        gameRepository,
		logger:                logger,
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
	regions, err := service.GameRepository.GetRegions()
	if err != nil {
		service.logger.Fatal("Failed to get game regions", slog.String("error", err.Error()))
	}

	var cronJobs []action.CronJob
	for _, region := range regions {
		regionID := region.ID
		time := region.ResetTime
		cronTime := fmt.Sprintf("%d %d * * *", time.Minute(), time.Hour())

		cronJob := action.NewCronJob(
			gocron.CronJob(cronTime, false),
			gocron.NewTask(service.AutoClaimTask, session, regionID),
		)

		cronJobs = append(cronJobs, cronJob)
	}

	return cronJobs
}

// Perform Genshin Impact daily check-in on HoYoLAB.
func (service *DailyService) DailyClaimCommandHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	discordUser := util.GetDiscordUser(interaction)
	discordID, _ := strconv.Atoi(discordUser.ID)

	user, err := service.HoyolabUserRepository.GetByDiscordID(discordID)
	if err != nil {
		content := "You are not registered yet, please register first."
		session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{Content: &content})
		service.logger.Error(content, slog.String("error", err.Error()))
		return
	}

	cookie := network.NewCookie(user.LtokenV2, user.LtmidV2, strconv.Itoa(user.ID))
	context := hoyolab.NewDailyRewardContext(hoyolab.Hk4eEndpoint, hoyolab.GenshinEventID, hoyolab.GenshinActID, hoyolab.GenshinSignGame)

	res, err := service.DailyRepository.Claim(cookie, context)
	if err != nil {
		content := "An internal error occurred while trying to check in. Please try again later."
		session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{Content: &content})
		service.logger.Error(content, slog.String("error", err.Error()))
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

// Task that automatically handles Genshin Impact daily check-in for all registered users.
func (service *DailyService) AutoClaimTask(session *discordgo.Session, regionID int) {
	offsetDiscordID := -1
	batchSize := 50

	failedUsers := []postgres.HoyolabUser{}

	for {
		users, err := service.HoyolabUserRepository.ListByRegionID(regionID, offsetDiscordID, batchSize)
		if err != nil {
			service.logger.Error("Failed to list users", slog.String("error", err.Error()))
			return
		}

		// If no more users are found, exit the loop. This means all users have been processed.
		if len(users) == 0 {
			break
		}

		failedUsers = append(failedUsers, service.autoClaim(session, users)...)

		// Start next batch from the last Discord ID in the current batch.
		offsetDiscordID = users[len(users)-1].DiscordID
	}

	retryCount := 0
	maxRetries := 3
	for len(failedUsers) > 0 {
		if retryCount >= maxRetries {
			break
		}

		failedUsers = service.autoClaim(session, failedUsers)
		retryCount++
	}

	if len(failedUsers) > 0 {
		discordIDs := make([]int, len(failedUsers))
		for i, user := range failedUsers {
			discordIDs[i] = user.DiscordID
		}

		service.logger.Error("Max retries reached for failed users", slog.Any("error", discordIDs))
		content := fmt.Sprintf(
			"%s\n%s",
			"We have encountered an error while trying to check in for you.",
			"Please re-register your account information using the `/register` command.",
		)

		for _, discordID := range discordIDs {
			service.sendChannelMessage(session, discordID, content)
		}
	}
}

// Automatically claim daily rewards for a list of users.
func (service *DailyService) autoClaim(session *discordgo.Session, users []postgres.HoyolabUser) []postgres.HoyolabUser {
	batchSize := 5
	failedUsers := []postgres.HoyolabUser{}

	for i, user := range users {
		// Process 5 users at a time before sleeping for 10 seconds to avoid rate limits.
		if i != 0 && i%batchSize == 0 {
			time.Sleep(time.Second * 10)
		}

		cookie := network.NewCookie(user.LtokenV2, user.LtmidV2, strconv.Itoa(user.ID))
		context := hoyolab.NewDailyRewardContext(hoyolab.Hk4eEndpoint, hoyolab.GenshinEventID, hoyolab.GenshinActID, hoyolab.GenshinSignGame)

		res, err := service.DailyRepository.Claim(cookie, context)

		content := "We have successfully checked in for you today!"
		if res.Retcode != 0 {
			content = res.Message
		} else if err != nil {
			service.logger.Warn(
				"Failed to auto claim daily reward",
				slog.Int("discordID", user.DiscordID),
				slog.String("error", err.Error()),
			)
			failedUsers = append(failedUsers, user)
			continue
		}

		service.sendChannelMessage(session, user.DiscordID, content)
	}

	return failedUsers
}

func (service *DailyService) sendChannelMessage(session *discordgo.Session, discordID int, content string) {
	channel, err := session.UserChannelCreate(strconv.Itoa(discordID))
	if err != nil {
		service.logger.Error(
			"Failed to send message to user channel",
			slog.Int("discordID", discordID),
			slog.String("error", err.Error()),
		)
		return
	}

	session.ChannelMessageSend(channel.ID, content)
}
