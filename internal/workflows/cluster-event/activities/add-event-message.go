package activities

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	v1 "k8s.io/api/core/v1"
	"log"
	"paul-go/internal"
)

const importantEventsChannelID = "931301531179966515"
const normalEventsChannelID = "931301737028001802"
const testEventsChannelID = "932115780768759878"

func AddEventMessage(ctx context.Context, event v1.Event) error {
	discordClient := internal.StartDiscord()
	defer discordClient.Close()
	discordClient.Identify.Intents = discordgo.IntentsGuildMessages
	err := discordClient.Open()
	if err != nil {
		log.Fatalln("Error opening Discord client connection,", err)
	}

	log.Println("Event received: ", event.ObjectMeta.Name)
	message := fmt.Sprintf("%v - %v %v\n", event.Type, event.InvolvedObject.Kind, event.Reason)
	message += fmt.Sprintf("```yaml\n")
	message += fmt.Sprintf("Namespace: %v\n", event.Namespace)
	message += fmt.Sprintf("Name: %v\n", event.Name)
	message += fmt.Sprintf("Message: %v\n", event.Message)
	message += fmt.Sprintf("```\n")
	message += fmt.Sprintf("||uid: `%v`||\n", event.UID)

	var destinationChannel string
	/*if event.Type == "Normal" {
		destinationChannel = normalEventsChannelID
	} else {
		destinationChannel = importantEventsChannelID
	}*/

	destinationChannel = testEventsChannelID

	_, sendError := discordClient.ChannelMessageSend(destinationChannel, message)
	if sendError != nil {
		log.Fatalln("Failed to send message: ", sendError)
	}

	return nil
}
