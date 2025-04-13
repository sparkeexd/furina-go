package util

import (
	"github.com/bwmarrin/discordgo"
)

// Get Discord user from Member or User.
// Member is only filled when the interaction is from a guild.
// User is only filled when the interaction is from a DM.
func GetDiscordUser(interaction *discordgo.InteractionCreate) *discordgo.User {
	if interaction.Member != nil {
		return interaction.Member.User
	}

	return interaction.User
}
