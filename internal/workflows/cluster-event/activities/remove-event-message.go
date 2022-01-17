package activities

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"log"
	"paul-go/internal"
)

func RemoveEventMessage(ctx context.Context, channelId string, eventMessageId string) error {
	discordClient := internal.StartDiscord()
	defer discordClient.Close()
	discordClient.Identify.Intents = discordgo.IntentsGuildMessages
	err := discordClient.Open()
	if err != nil {
		log.Println("Error opening Discord client connection,", err)
		return err
	}

	log.Println("deleting: message")
	deleteErr := discordClient.ChannelMessageDelete(channelId, eventMessageId)
	return deleteErr
}
