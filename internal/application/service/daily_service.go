package service

import (
	"fmt"
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
var DailyClaimContext = map[string]hoyolab.DailyRewardContext{
	postgres.GenshinEnum: hoyolab.NewDailyRewardContext(
		hoyolab.Hk4eEndpoint,
		hoyolab.GenshinEventID,
		hoyolab.GenshinActID,
		hoyolab.GenshinSignGame,
	),
	postgres.StarRailEnum: hoyolab.NewDailyRewardContext(
		hoyolab.SgPublicEndpoint,
		hoyolab.StarRailEventID,
		hoyolab.StarRailActID,
		hoyolab.StarRailSignGame,
	),
	postgres.ZenlessEnum: hoyolab.NewDailyRewardContext(
		hoyolab.SgPublicEndpoint,
		hoyolab.ZenlessEventID,
		hoyolab.ZenlessActID,
		hoyolab.ZenlessSignGame,
	),
}

// Service that handles daily check-in commands.
type DailyService struct {
	dailyRepository       hoyolab.DailyRepository
	userRepository        postgres.UserRepository
	accountRepository     postgres.AccountRepository
	interactionRepository discord.InteractionRepository
	logger                *logger.Logger
}

// Create a new daily service.
func NewDailyService(
	dailyRepository hoyolab.DailyRepository,
	userRepository postgres.UserRepository,
	accountRepository postgres.AccountRepository,
	interactionRepository discord.InteractionRepository,
	logger *logger.Logger,
) DailyService {
	return DailyService{
		dailyRepository:       dailyRepository,
		userRepository:        userRepository,
		accountRepository:     accountRepository,
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
	cronTab := "0 20 * * *"
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

	embed := service.dailyClaim(discordID)

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
		users, err := service.userRepository.ListDiscordUsers(offsetDiscordID, batchSize)
		if err != nil {
			service.logger.Error("Failed to list users", slog.String("error", err.Error()))
			return
		}

		// If no more users are found, exit the loop. This means all users have been processed.
		if len(users) == 0 {
			break
		}

		for _, user := range users {
			// Sleep for 10 seconds per user to mitigate rate limits.
			time.Sleep(time.Second * 10)
			embed := service.dailyClaim(user.ID)
			service.sendChannelMessageEmbed(session, user.ID, embed.MessageEmbed)
		}

		// Start next batch from the last Discord ID in the current batch.
		offsetDiscordID = users[len(users)-1].ID
	}
}

// Claim daily rewards for a user.
func (service *DailyService) dailyClaim(discordID int) (embed *entity.Embed) {
	service.logger.Info("Claiming daily reward", slog.Int("discordID", discordID))

	user, accounts, err := service.getUserAccounts(discordID)
	if err != nil {
		description := "Your game accounts could not be fetched. Please try registering again."
		embed = service.interactionRepository.CreateErrorEmbed().SetDescription(description)
		return embed
	}

	cookie := network.NewCookie(user.LtokenV2, user.LtmidV2, user.LtuidV2)
	fields := []*discordgo.MessageEmbedField{}

	for _, account := range accounts {
		context, exists := DailyClaimContext[account.Game]
		if !exists {
			service.logger.Error(
				"Failed to fetch daily claim context: Invalid Game ID",
				slog.String("gameID", account.Game),
			)
			return
		}

		res, err := service.dailyRepository.Claim(cookie, context)

		gameTitle := service.accountRepository.GetGameTitle(account.Game)
		content := service.getEmbedContent(res, err, discordID)
		fields = append(fields, &discordgo.MessageEmbedField{Name: gameTitle, Value: content})
	}

	value := ""
	for _, field := range fields {
		value = fmt.Sprintf("%s\n**%s**: %s", value, field.Name, field.Value)
	}

	embed = service.interactionRepository.CreateEmbed().
		SetTitle("Daily Check-in").
		SetDescription(fmt.Sprintf("Claim your HoYoLAB daily check-in rewards!\n%s", value)).
		SetThumbnail("https://media.tenor.com/EhXA2CCJ-QUAAAAj/furina.gif")

	return embed
}

// Get user and their HoYoverse game accounts.
func (service *DailyService) getUserAccounts(discordID int) (user entity.User, accounts []entity.Account, err error) {
	user, err = service.userRepository.GetByDiscordID(discordID)
	if err != nil {
		service.logger.Error(
			"Failed to get user",
			slog.Int("discordID", discordID),
			slog.String("error", err.Error()),
		)

		return user, accounts, err
	}

	accounts, err = service.accountRepository.ListByDiscordID(discordID)
	if err != nil {
		service.logger.Error(
			"Failed to get user's HoYoverse game accounts",
			slog.Int("discordID", discordID),
			slog.String("error", err.Error()),
		)

		return user, accounts, err
	}

	return user, accounts, err
}

// Get embed content.
func (service *DailyService) getEmbedContent(res entity.DailyClaim, err error, discordID int) (content string) {
	content = ""
	if err != nil {
		content = err.Error()
		service.logger.Warn(
			"Failed to auto claim daily reward",
			slog.Int("discordID", discordID),
			slog.String("error", err.Error()),
		)
	} else {
		switch res.Retcode {
		case hoyolab.OK, hoyolab.DailyAlreadyClaimed:
			content = "Reward claimed!"
		case hoyolab.InvalidCookie:
			content = "Cookie is invalid/has expired."
		default:
			content = res.Message
		}

		service.logger.Info(
			res.Message,
			slog.Int("discordID", discordID),
			slog.Int("retcode", res.Retcode),
		)
	}

	return content
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
