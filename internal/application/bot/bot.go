package bot

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sparkeexd/mimo/internal/application/service"
	"github.com/sparkeexd/mimo/internal/domain/action"
	"github.com/sparkeexd/mimo/internal/infrastructure/hoyolab"
	"github.com/sparkeexd/mimo/internal/infrastructure/postgres"
)

// Discord bot.
type Bot struct {
	Token           string
	Session         *discordgo.Session
	CommandServices []action.CommandService
	JobServices     []action.JobService
	Scheduler       gocron.Scheduler
}

// Create a new Discord bot.
func NewBot() Bot {
	token := os.Getenv("BOT_TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to initialize scheduler: %v", err)
	}

	dailyRepository := hoyolab.NewDailyRepository()
	userRepository := postgres.NewHoyolabUserRepository(db)

	pingService := service.NewPingService()
	dailyService := service.NewDailyService(dailyRepository, userRepository)

	commandServices := []action.CommandService{&pingService, &dailyService}
	jobServices := []action.JobService{&dailyService}

	bot := Bot{
		Token:           token,
		Session:         session,
		CommandServices: commandServices,
		JobServices:     jobServices,
		Scheduler:       scheduler,
	}

	return bot
}

// Start Discord bot.
func (bot *Bot) Start() {
	log.Println("Creating discord bot session...")
	err := bot.Session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Registering commands...")
	bot.registerCommands()

	log.Println("Registering jobs...")
	bot.registerJobs()
	bot.Scheduler.Start()

	bot.Session.AddHandler(bot.Ready)

	// Event listener to stop the bot.
	log.Println("Bot is now running! Press Ctrl+C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Closing discord bot session...")
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
		jobs := service.Jobs(bot.Session)

		for _, job := range jobs {
			bot.Scheduler.NewJob(job.Definition, job.Task)
		}
	}
}
