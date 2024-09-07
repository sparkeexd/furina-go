package bot

import (
	"flag"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/hoyoapi/client"
	"github.com/sparkeexd/mimo/commands"
)

var (
	// Discord bot Session.
	Session *discordgo.Session

	// Discord bot parameters.
	BotToken = flag.String("token", "", "Bot access token.")

	// Hoyo API clients.
	GenshinClient  *client.GenshinClient
	StarRailClient *client.StarRailClient
	ZenlessClient  *client.ZenlessClient
)

// Create discord bot session.
func CreateSession() {
	flag.Parse()

	var err error
	Session, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}

	Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err = Session.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}
}

// Initializing Hoyo API clients.
func InitializeClients() {
	language := "en-us"
	clientOptions := client.NewClientOptions().
		AddLanguage(language).
		Build()

	GenshinClient = client.NewGenshinClient(clientOptions)
	StarRailClient = client.NewStarRailClient(clientOptions)
	ZenlessClient = client.NewZenlessClient(clientOptions)
}

// Register the slash command.
// Add a handler executes the registered handler if its corresponding command exists.
func AddCommands(commands map[string]commands.Command) {
	Session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if command, exists := commands[interaction.ApplicationCommandData().Name]; exists {
			command.Handler(session, interaction)
		}
	})

	for _, v := range commands {
		_, err := Session.ApplicationCommandCreate(Session.State.User.ID, "", v.Command)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Command.Name, err)
		}
	}
}

// Close the session.
func CloseSession() {
	Session.Close()
}
