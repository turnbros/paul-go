package activities

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"paul-go/internal"
	"paul-go/internal/util"
)

func AddEventMessage(ctx context.Context, channelId string, event util.ClusterEventMessage) (string, error) {
	discordClient := internal.StartDiscord()
	defer discordClient.Close()
	discordClient.Identify.Intents = discordgo.IntentsGuildMessages
	err := discordClient.Open()
	if err != nil {
		log.Println("Error opening Discord client connection,", err)
		return "", err
	}

	log.Println("Event received: ", event.EventName)
	message := fmt.Sprintf(":information_source: %v - %v %v\n", event.EventType, event.ObjectKind, event.EventReason)
	message += fmt.Sprintf("||event uid: `%v`||\n", event.EventUID)
	message += fmt.Sprintf("```yaml\n")
	message += fmt.Sprintf("Namespace: %v\n", event.ObjectNamespace)
	message += fmt.Sprintf("Name: %v\n", event.ObjectName)
	message += fmt.Sprintf("Message: %v\n", event.EventMessage)
	message += fmt.Sprintf("Count: %v\n", event.EventCount)
	message += fmt.Sprintf("```\n")
	message += fmt.Sprintf("-\n")

	log.Println(message)
	eventMessage, sendError := discordClient.ChannelMessageSend(channelId, message)
	if sendError != nil {
		log.Println("Failed to send message: ", sendError)
		return "", sendError
	}
	return eventMessage.ID, nil
}
