package service

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"github.com/sparkeexd/mimo/internal/domain/entity"
	"github.com/sparkeexd/mimo/internal/infrastructure/discord"
	"github.com/sparkeexd/mimo/internal/infrastructure/hoyolab"
	"github.com/sparkeexd/mimo/internal/infrastructure/postgres"
	"github.com/sparkeexd/mimo/pkg/logger"
	"github.com/sparkeexd/mimo/pkg/network"
)

// Mapping HoYoverse game to its respective daily reward context.
var DailyClaimContext = map[int]hoyolab.DailyRewardContext{
	1: hoyolab.NewDailyRewardContext(hoyolab.Hk4eEndpoint, hoyolab.GenshinEventID, hoyolab.GenshinActID, hoyolab.GenshinSignGame),
	2: hoyolab.NewDailyRewardContext(hoyolab.SgPublicEndpoint, hoyolab.StarRailEventID, hoyolab.StarRailActID, hoyolab.StarRailSignGame),
	3: hoyolab.NewDailyRewardContext(hoyolab.SgPublicEndpoint, hoyolab.ZenlessEventID, hoyolab.ZenlessActID, hoyolab.ZenlessSignGame),
}

// Service that handles daily check-in commands.
type DailyService struct {
	dailyRepository       hoyolab.DailyRepository
	discordUserRepository postgres.DiscordUserRepository
	hoyolabUserRepository postgres.HoyolabUserRepository
	gameUserRepository    postgres.GameUserRepository
	interactionRepository discord.InteractionRepository
	logger                *logger.Logger
}

// Create a new daily service.
func NewDailyService(
	dailyRepository hoyolab.DailyRepository,
	discordUserRepository postgres.DiscordUserRepository,
	hoyolabUserRepository postgres.HoyolabUserRepository,
	gameRepository postgres.GameUserRepository,
	interactionRepository discord.InteractionRepository,
	logger *logger.Logger,
) DailyService {
	return DailyService{
		dailyRepository:       dailyRepository,
		discordUserRepository: discordUserRepository,
		hoyolabUserRepository: hoyolabUserRepository,
		gameUserRepository:    gameRepository,
		interactionRepository: interactionRepository,
		logger:                logger,
	}
}

// Service's slash commands to be registered.
func (service *DailyService) Commands() map[string]entity.Command {
	return map[string]entity.Command{
		"daily": entity.NewCommand(
			&discordgo.ApplicationCommand{
				Name:        "daily",
				Description: "Claim your HoYoLAB daily check-in rewards.",
			},
			service.dailyClaimCommandHandler,
		),
	}
}

// Service's cron jobs to be registered.
func (service *DailyService) Jobs(session *discordgo.Session) []entity.CronJob {
	cronJobs := []entity.CronJob{}
	jobName := "AutoDailyClaim"
	cronTab := "0 16 * * *"
	cronJob := entity.NewCronJob(
		gocron.CronJob(cronTab, false),
		gocron.NewTask(service.autoDailyClaimTask, session),
		gocron.WithName(jobName),
		cronTab,
	)

	cronJobs = append(cronJobs, cronJob)

	return cronJobs
}

// Perform Genshin Impact daily check-in on HoYoLAB.
func (service *DailyService) dailyClaimCommandHandler(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})

	user := service.interactionRepository.GetDiscordUser(interaction)
	discordID, _ := strconv.Atoi(user.ID)

	embed, err := service.dailyClaim(discordID)
	if err != nil {
		// Send error embed
		return
	}

	// res, err := service.dailyRepository.Claim(cookie, context)
	// if err != nil {
	// 	content := "An internal error occurred while trying to check in. Please try again later."
	// 	session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{Content: &content})
	// 	service.logger.Error(content, slog.String("error", err.Error()))
	// 	return
	// }

	// content := fmt.Sprintf("You have successfully checked in, %s!", user.Mention())
	// if res.Retcode != 0 {
	// 	content = res.Message
	// 	service.logger.Info(
	// 		"Failed to auto claim daily reward",
	// 		slog.Int("discordID", discordID),
	// 		slog.Int("retcode", res.Retcode),
	// 		slog.String("message", res.Message),
	// 	)
	// }

	session.InteractionResponseEdit(interaction.Interaction, &discordgo.WebhookEdit{
		Embeds: &[]*discordgo.MessageEmbed{embed.MessageEmbed},
	})
}

// Task that automatically handles Genshin Impact daily check-in for all registered users.
func (service *DailyService) autoDailyClaimTask(session *discordgo.Session) {
	service.logger.Info("Running auto claim task")

	offsetDiscordID := -1
	batchSize := 50

	for {
		discordUsers, err := service.discordUserRepository.ListDiscordUsers(offsetDiscordID, batchSize)
		if err != nil {
			service.logger.Error("Failed to list users", slog.String("error", err.Error()))
			return
		}

		// If no more users are found, exit the loop. This means all users have been processed.
		if len(discordUsers) == 0 {
			break
		}

		for _, discordUser := range discordUsers {
			// Sleep for 10 seconds per user to mitigate rate limits.
			time.Sleep(time.Second * 10)

			embed, err := service.dailyClaim(discordUser.ID)
			if err != nil {
				// Send error embed
				continue
			}

			service.sendChannelMessageEmbed(session, discordUser.ID, embed.MessageEmbed)
		}

		// Start next batch from the last Discord ID in the current batch.
		offsetDiscordID = discordUsers[len(discordUsers)-1].ID
	}
}

// Claim daily rewards for a Discord user.
func (service *DailyService) dailyClaim(discordID int) (*entity.Embed, error) {
	service.logger.Info("Auto claiming daily reward", slog.Int("discordID", discordID))

	hoyolabUser, err := service.hoyolabUserRepository.GetByDiscordID(discordID)
	if err != nil {
		service.logger.Error(
			"Failed to get Discord user's HoYoLAB account",
			slog.Int("discordID", discordID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	gameUsers, err := service.gameUserRepository.ListByDiscordID(discordID)
	if err != nil {
		service.logger.Error(
			"Failed to get Discord user's HoYoverse game accounts",
			slog.Int("discordID", discordID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	cookie := network.NewCookie(hoyolabUser.LtokenV2, hoyolabUser.LtmidV2, strconv.Itoa(hoyolabUser.ID))
	fields := []*discordgo.MessageEmbedField{{Name: "", Value: ""}}

	for _, gameUser := range gameUsers {
		context := DailyClaimContext[gameUser.GameID]
		res, err := service.dailyRepository.Claim(cookie, context)

		content := "Successfully checked in."
		if res.Retcode != 0 {
			content = res.Message
			service.logger.Info(
				"Failed to auto claim daily reward",
				slog.Int("discordID", discordID),
				slog.Int("retcode", res.Retcode),
				slog.String("message", res.Message),
			)
		} else if err != nil {
			content = err.Error()
			service.logger.Warn(
				"Failed to auto claim daily reward",
				slog.Int("discordID", discordID),
				slog.String("error", err.Error()),
			)
		}

		game, _ := service.gameUserRepository.GetGameByID(gameUser.GameID)
		fields = append(
			fields,
			&discordgo.MessageEmbedField{Name: game.Name, Value: content},
			&discordgo.MessageEmbedField{Name: "", Value: ""}, // To add small newline between fields
		)
	}

	embed := service.interactionRepository.CreateEmbed().
		SetTitle("Daily Check-in").
		SetDescription("Claim your HoYoLAB daily check-in rewards!").
		SetThumbnail("https://media.tenor.com/EhXA2CCJ-QUAAAAj/furina.gif").
		AddFields(fields)

	return embed, nil
}

// Send message to user channel.
func (service *DailyService) sendChannelMessageEmbed(session *discordgo.Session, discordID int, embed *discordgo.MessageEmbed) error {
	channel, err := session.UserChannelCreate(strconv.Itoa(discordID))
	if err != nil {
		service.logger.Error(
			"Failed to send message to user channel",
			slog.Int("discordID", discordID),
			slog.String("error", err.Error()),
		)
		return err
	}

	_, err = session.ChannelMessageSendEmbed(channel.ID, embed)
	return err
}
