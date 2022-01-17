package activities

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"paul-go/internal"
	"paul-go/internal/util"
)

func UpdateEventMessage(ctx context.Context, channelId string, eventMessageId string, event util.ClusterEventMessage) error {
	discordClient := internal.StartDiscord()
	defer discordClient.Close()
	discordClient.Identify.Intents = discordgo.IntentsGuildMessages
	err := discordClient.Open()
	if err != nil {
		log.Println("Error opening Discord client connection,", err)
		return err
	}

	var messageIcon string
	if event.EventType == "Normal" {
		messageIcon = ":information_source:"
	} else {
		messageIcon = ":warning:"
	}

	log.Println("Event received: ", event.EventName)
	message := fmt.Sprintf("%v %v - %v %v\n", messageIcon, event.EventType, event.ObjectKind, event.EventReason)
	message += fmt.Sprintf("||event uid: `%v`||\n", event.EventUID)
	message += fmt.Sprintf("```yaml\n")
	message += fmt.Sprintf("Name: %v\n", event.ObjectName)
	message += fmt.Sprintf("Namespace: %v\n", event.ObjectNamespace)
	message += fmt.Sprintf("Message: \"%v\"\n", event.EventMessage)
	message += fmt.Sprintf("Count: %v\n", event.EventCount)
	message += fmt.Sprintf("```\n")
	message += fmt.Sprintf("-\n")

	_, sendError := discordClient.ChannelMessageEdit(channelId, eventMessageId, message)
	if sendError != nil {
		log.Println("Failed to send message: ", sendError)
		return sendError
	}
	return nil
}
