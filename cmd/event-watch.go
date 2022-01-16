package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"paul-go/internal"
	"paul-go/internal/util"
)

const importantEventsChannelID = "931301531179966515"
const normalEventsChannelID = "931301737028001802"

func main() {
	discordClient := internal.StartDiscord()
	defer discordClient.Close()
	discordClient.Identify.Intents = discordgo.IntentsGuildMessages
	err := discordClient.Open()
	if err != nil {
		log.Fatalln("Error opening Discord client connection,", err)
	}

	kubeClient := util.GetKubeClient()
	ctx := context.Background()

	eventList, err := kubeClient.CoreV1().Events(metav1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to list events: ", err)
	}

	eventWatcher, err := kubeClient.CoreV1().Events(metav1.NamespaceAll).Watch(ctx, metav1.ListOptions{ResourceVersion: eventList.ListMeta.ResourceVersion})
	if err != nil {
		log.Fatalln("Failed to start the event watcher: ", err)
	}

	log.Println("Waiting for cluster events...")
	for event := range eventWatcher.ResultChan() {
		ev := event.Object.(*v1.Event)
		switch event.Type {
		case watch.Added:
			log.Println("Event received: ", ev.ObjectMeta.Name)
			message := fmt.Sprintf("%v - %v %v\n", ev.Type, ev.InvolvedObject.Kind, ev.Reason)
			message += fmt.Sprintf("```yaml\n")
			message += fmt.Sprintf("Namespace: %v\n", ev.Namespace)
			message += fmt.Sprintf("Name: %v\n", ev.Name)
			message += fmt.Sprintf("Message: %v\n", ev.Message)
			message += fmt.Sprintf("```\n")

			var destinationChannel string
			if ev.Type == "Normal" {
				destinationChannel = normalEventsChannelID
			} else {
				destinationChannel = importantEventsChannelID
			}

			_, sendError := discordClient.ChannelMessageSend(destinationChannel, message)
			if sendError != nil {
				log.Fatalln("Failed to send message: ", sendError)
			}
		case watch.Error:
			log.Printf("Error Event %s \n", ev.ObjectMeta.Name)
		case watch.Bookmark:
			log.Printf("Bookmark Event %s \n", ev.ObjectMeta.Name)
		case watch.Modified:
			log.Printf("Modified Event %s \n", ev.ObjectMeta.Name)
			log.Printf("Modified Event %s \n", ev.Type)
		case watch.Deleted:
			log.Printf("Deleted Event %s \n", ev.ObjectMeta.Name)
			log.Printf("Modified Event %s \n", ev.Type)
		}
	}
}
