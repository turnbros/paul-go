package main

import (
	"github.com/bwmarrin/discordgo"
	"go.temporal.io/sdk/client"
	"log"
	"os"
	"os/signal"
	"paul/internal"
	"strings"
	"syscall"
)

var discordClient *discordgo.Session
var temporalClient client.Client

const dialogflowProjectId = "paul-fmma"
const dialogflowLanguageCode = "en"

func main() {

	temporalClient = internal.StartTemporal()
	defer temporalClient.Close()

	discordClient = internal.StartDiscord()
	defer discordClient.Close()

	discordClient.AddHandler(HandleMessage)
	discordClient.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err := discordClient.Open()
	if err != nil {
		log.Fatalln("error opening Discord client connection,", err)
		panic(5)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	botAuthorId := "<@!" + s.State.User.ID + ">"
	message := strings.ToLower(m.Content)
	if strings.HasPrefix(message, botAuthorId) {
		requestMessage := strings.Replace(message, botAuthorId, "", 1)
		intentAction, intentParameters, paulResponse := internal.ParseRequest(dialogflowProjectId, requestMessage, dialogflowLanguageCode)

		if strings.HasPrefix(intentAction, "workflow") {
			paulResponse = internal.ExecuteWorkflow(temporalClient, intentAction, intentParameters)
		}

		s.ChannelMessageSend(m.ChannelID, paulResponse)
	}
}
