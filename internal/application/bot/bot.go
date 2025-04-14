package bot

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sparkeexd/mimo/internal/application/service"
	"github.com/sparkeexd/mimo/internal/application/util"
	"github.com/sparkeexd/mimo/internal/domain/action"
	"github.com/sparkeexd/mimo/internal/domain/logger"
	"github.com/sparkeexd/mimo/internal/infrastructure/hoyolab"
	"github.com/sparkeexd/mimo/internal/infrastructure/postgres"
)

// Discord bot.
type Bot struct {
	Session         *discordgo.Session
	CommandServices []action.CommandService
	JobServices     []action.JobService
	Scheduler       gocron.Scheduler
	Logger          *logger.Logger
}

// Create a new Discord bot.
func NewBot() Bot {
	logger := logger.NewLogger()

	token := os.Getenv("BOT_TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Fatal("Invalid bot parameters", slog.String("error", err.Error()))
	}

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("Unable to connect to database", slog.String("error", err.Error()))
	}

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		logger.Fatal("Failed to initialize scheduler", slog.String("error", err.Error()))
	}

	gameRepository := postgres.NewGameRepository(db)
	dailyRepository := hoyolab.NewDailyRepository(logger)
	userRepository := postgres.NewHoyolabUserRepository(db)

	pingService := service.NewPingService()
	dailyService := service.NewDailyService(dailyRepository, userRepository, gameRepository, logger)

	commandServices := []action.CommandService{&pingService, &dailyService}
	jobServices := []action.JobService{&dailyService}

	bot := Bot{
		Session:         session,
		CommandServices: commandServices,
		JobServices:     jobServices,
		Scheduler:       scheduler,
		Logger:          logger,
	}

	return bot
}

// Start Discord bot.
func (bot *Bot) Start() {
	bot.Session.AddHandler(bot.Ready)

	bot.Logger.Info("Creating discord bot session...")
	err := bot.Session.Open()
	if err != nil {
		bot.Logger.Fatal("Cannot open the session", slog.String("error", err.Error()))
	}

	bot.Logger.Info("Registering commands...")
	bot.registerCommands()

	bot.Logger.Info("Registering jobs...")
	bot.registerJobs()
	bot.Scheduler.Start()

	// Event listener to stop the bot.
	bot.Logger.Info("Bot is now running! Press Ctrl+C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	bot.Logger.Info("Closing discord bot session...")
	bot.Session.Close()
	bot.Scheduler.Shutdown()
}

// Register the slash commands.
// Middleware is attached to each command to block interactions outside of the guild.
// Requires reloading Discord client to view the changes.
func (bot *Bot) registerCommands() {
	var commandsToRegister []*discordgo.ApplicationCommand

	for _, service := range bot.CommandServices {
		commands := service.Commands()

		bot.Session.AddHandler(
			func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
				if command, exists := commands[interaction.ApplicationCommandData().Name]; exists {
					discordUser := util.GetDiscordUser(interaction)
					bot.Logger.Info(
						"Command invoked",
						slog.String("command", command.Command.Name),
						slog.Group("user",
							slog.String("id", discordUser.ID),
							slog.String("name", discordUser.Username),
						),
						slog.Group("guild",
							slog.String("id", interaction.GuildID),
							slog.String("name", interaction.ChannelID),
						),
					)

					bot.InteractionCreate(command.Handler)(session, interaction)
				}
			},
		)

		for _, command := range commands {
			commandsToRegister = append(commandsToRegister, command.Command)
		}
	}

	// Overwrite all existing commands, which allow clearing out old commands.
	bot.Session.ApplicationCommandBulkOverwrite(bot.Session.State.User.ID, "", commandsToRegister)
}

// Register the cron jobs.
func (bot *Bot) registerJobs() {
	for _, service := range bot.JobServices {
		cronJobs := service.Jobs(bot.Session)

		for _, cronJob := range cronJobs {
			job, err := bot.Scheduler.NewJob(cronJob.Definition, cronJob.Task, cronJob.Option)
			if err != nil {
				bot.Logger.Error("Failed to register cron job", slog.String("error", err.Error()))
			}

			bot.Logger.Info(
				"Registered cron job",
				slog.String("name", job.Name()),
				slog.String("crontab", cronJob.CronTab),
			)

		}
	}
}
