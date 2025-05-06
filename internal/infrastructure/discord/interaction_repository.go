package discord

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sparkeexd/furina/internal/domain/entity"
	"github.com/sparkeexd/furina/pkg/logger"
)

// Repository for handling Discord interactions.
type InteractionRepository struct {
	logger *logger.Logger
}

// Create a new Discord interaction repository.
func NewInteractionRepository(logger *logger.Logger) InteractionRepository {
	return InteractionRepository{
		logger: logger,
	}
}

// Get Discord user from Member or User.
// Member is only filled when the interaction is from a guild.
// User is only filled when the interaction is from a DM.
func (repository InteractionRepository) GetDiscordUser(interaction *discordgo.InteractionCreate) *discordgo.User {
	if interaction.Member != nil {
		return interaction.Member.User
	}

	return interaction.User
}

// Create a user channel to send a DM message to them.
func (repository InteractionRepository) CreateUserChannel(session *discordgo.Session, discordID int) (*discordgo.Channel, error) {
	channel, err := session.UserChannelCreate(strconv.Itoa(discordID))
	if err != nil {
		repository.logger.Error(
			"Failed to send message to user channel",
			slog.Int("discordID", discordID),
			slog.String("error", err.Error()),
		)
		return nil, err
	}

	return channel, nil
}

// Create a normal embed.
func (repository InteractionRepository) CreateEmbed() *entity.Embed {
	embed := entity.NewEmbed().
		SetTimestamp(time.Now().Format(time.RFC3339)).
		SetColor(0x1A1F68)

	return embed
}

// Create an error embed.
func (repository InteractionRepository) CreateErrorEmbed() *entity.Embed {
	embed := entity.NewEmbed().
		SetTitle("Error").
		SetThumbnail("https://media.tenor.com/7nsbCGbleT0AAAAi/furina-genshin-impact.png").
		SetColor(0xFC2C2C)

	return embed
}
