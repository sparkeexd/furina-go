package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron/v2"
	"github.com/sparkeexd/mimo/internal/database"
	"github.com/sparkeexd/mimo/internal/middleware"
	"github.com/sparkeexd/mimo/internal/models"
)

// Discord bot.
type Bot struct {
	Token     string
	Session   *discordgo.Session
	DB        *database.DB
	Commands  []map[string]models.Command
	Jobs      []models.CronJob
	Scheduler gocron.Scheduler
}

// Create a new Discord bot.
func NewBot() *Bot {
	token := os.Getenv("BOT_TOKEN")
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	db, err := database.DatabaseClient()
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	scheduler, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Failed to initialize scheduler: %v", err)
	}

	bot := &Bot{
		Token:     token,
		Session:   session,
		Scheduler: scheduler,
		DB:        db,
	}

	bot.Session.AddHandler(middleware.Ready)

	return bot
}

// Start Discord bot.
func (bot *Bot) Start(commands []map[string]models.Command, jobs []models.CronJob) {
	log.Println("Creating discord bot session...")
	err := bot.Session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Registering commands...")
	bot.RegisterCommands(commands)

	log.Println("Registering jobs...")
	bot.RegisterJobs(jobs)
	bot.Scheduler.Start()

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
func (bot *Bot) RegisterCommands(commands []map[string]models.Command) {
	var commandsToRegister []*discordgo.ApplicationCommand

	for _, commands := range commands {
		bot.Session.AddHandler(
			func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
				if command, exists := commands[interaction.ApplicationCommandData().Name]; exists {
					middleware.InteractionCreate(command.Handler)(session, interaction, bot.DB)
				}
			},
		)

		for _, v := range commands {
			commandsToRegister = append(commandsToRegister, v.Command)
		}
	}

	// Overwrite all existing commands, which allow clearing out old commands.
	bot.Session.ApplicationCommandBulkOverwrite(bot.Session.State.User.ID, "", commandsToRegister)
}

// Register the cron jobs.
func (bot *Bot) RegisterJobs(jobs []models.CronJob) {
	for _, job := range jobs {
		bot.Scheduler.NewJob(job.Definition, job.Task)
	}
}
