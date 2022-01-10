package internal

import (
	"github.com/bwmarrin/discordgo"
	"os"
)

func StartDiscord() *discordgo.Session {
	discordAuthToken := os.Getenv("DISCORD_TOKEN")
	discord, err := discordgo.New("Bot " + discordAuthToken)
	if err != nil {
		panic(err)
	}
	return discord
}
