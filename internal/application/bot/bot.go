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
	"github.com/sparkeexd/mimo/internal/domain/abstract"
	"github.com/sparkeexd/mimo/internal/infrastructure/discord"
	"github.com/sparkeexd/mimo/internal/infrastructure/hoyolab"
	"github.com/sparkeexd/mimo/internal/infrastructure/postgres"
	"github.com/sparkeexd/mimo/pkg/logger"
)

// Discord bot.
type Bot struct {
	session         *discordgo.Session
	commandServices []abstract.CommandService
	jobServices     []abstract.JobService
	scheduler       gocron.Scheduler
	logger          *logger.Logger
}

// Create a new Discord bot.
func NewBot() Bot {
	context := context.Background()
	logger := logger.NewLogger()

	token := os.Getenv("BOT_TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Fatal("Invalid bot parameters", slog.String("error", err.Error()))
	}

	db, err := pgxpool.New(context, os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("Unable to connect to database", slog.String("error", err.Error()))
	}

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		logger.Fatal("Failed to initialize scheduler", slog.String("error", err.Error()))
	}

	commandServices, jobServices := initializeServices(db, logger)

	bot := Bot{
		session:         session,
		commandServices: commandServices,
		jobServices:     jobServices,
		scheduler:       scheduler,
		logger:          logger,
	}

	return bot
}

// Initialize services used by the Discord bot.
func initializeServices(db *pgxpool.Pool, logger *logger.Logger) ([]abstract.CommandService, []abstract.JobService) {
	dailyRepository := hoyolab.NewDailyRepository(logger)
	userRepository := postgres.NewUserRepository(db)
	accountRepository := postgres.NewAccountRepository(db)
	interactionRepository := discord.NewInteractionRepository(logger)

	pingService := service.NewPingService(interactionRepository)
	dailyService := service.NewDailyService(
		dailyRepository,
		userRepository,
		accountRepository,
		interactionRepository,
		logger,
	)

	commandServices := []abstract.CommandService{&pingService, &dailyService}
	jobServices := []abstract.JobService{&dailyService}

	return commandServices, jobServices
}

// Start Discord bot.
func (bot *Bot) Start() {
	bot.session.AddHandler(bot.logReady)

	bot.logger.Info("Creating discord bot session...")
	err := bot.session.Open()
	if err != nil {
		bot.logger.Fatal("Cannot open the session", slog.String("error", err.Error()))
	}

	bot.logger.Info("Registering commands...")
	bot.registerCommands()

	bot.logger.Info("Registering jobs...")
	bot.registerJobs()
	bot.scheduler.Start()

	// Event listener to stop the bot.
	bot.logger.Info("Bot is now running! Press Ctrl+C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	bot.logger.Info("Closing discord bot session...")
	bot.session.Close()
	bot.scheduler.Shutdown()
}

// Register the slash commands.
// Middleware is attached to each command to block interactions outside of the guild.
// Requires reloading Discord client to view the changes.
func (bot *Bot) registerCommands() {
	var commandsToRegister []*discordgo.ApplicationCommand

	for _, service := range bot.commandServices {
		commands := service.Commands()

		bot.session.AddHandler(
			func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
				if command, exists := commands[interaction.ApplicationCommandData().Name]; exists {
					user := interaction.User
					if user == nil {
						user = interaction.Member.User
					}

					bot.logger.Info(
						"Command invoked",
						slog.String("command", command.Command.Name),
						slog.Group("user",
							slog.String("id", user.ID),
							slog.String("name", user.Username),
						),
						slog.Group("guild",
							slog.String("id", interaction.GuildID),
							slog.String("name", interaction.ChannelID),
						),
					)

					bot.filterInteraction(command.Handler, user)(session, interaction)
				}
			},
		)

		for _, command := range commands {
			commandsToRegister = append(commandsToRegister, command.Command)
		}
	}

	// Overwrite all existing commands, which allow clearing out old commands.
	bot.session.ApplicationCommandBulkOverwrite(bot.session.State.User.ID, "", commandsToRegister)
}

// Register the cron jobs.
func (bot *Bot) registerJobs() {
	for _, service := range bot.jobServices {
		cronJobs := service.Jobs(bot.session)

		for _, cronJob := range cronJobs {
			job, err := bot.scheduler.NewJob(cronJob.Definition, cronJob.Task, cronJob.Option)
			if err != nil {
				bot.logger.Error("Failed to register cron job", slog.String("error", err.Error()))
			}

			bot.logger.Info(
				"Registered cron job",
				slog.Group("job",
					slog.String("name", job.Name()),
					slog.String("crontab", cronJob.CronTab),
				),
			)

		}
	}
}
