package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"go.temporal.io/sdk/client"
	"log"
	"os"
	"os/signal"
	"paul/internal"
	"paul/internal/dialogflow"
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

type MessageLogEntry struct {
	Timestamp      string `json:"timestamp"`
	MsgID          string `json:"msg_id"`
	MsgChannelID   string `json:"msg_channel_id"`
	MsgGuildID     string `json:"msg_guild_id"`
	MsgContent     string `json:"msg_content"`
	AuthorID       string `json:"author_id"`
	AuthorBot      bool   `json:"author_bot"`
	AuthorEmail    string `json:"author_email"`
	AuthorUsername string `json:"author_username"`
}

func logIncomingMessage(message *discordgo.MessageCreate) {
	// Create a message log entry
	messageLogEntry := MessageLogEntry{
		Timestamp:      string(message.Timestamp),
		MsgID:          message.ID,
		MsgChannelID:   message.ChannelID,
		MsgGuildID:     message.GuildID,
		MsgContent:     message.Content,
		AuthorID:       message.Author.ID,
		AuthorBot:      message.Author.Bot,
		AuthorEmail:    message.Author.Email,
		AuthorUsername: message.Author.Username,
	}
	messageLogEntryJson, err := json.Marshal(messageLogEntry)
	if err != nil {
		log.Fatalln("Failed to marshall the message log: ", err)
	}
	log.Println(string(messageLogEntryJson))
}

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	logIncomingMessage(m)
	if m.Author.ID == s.State.User.ID {
		return
	}

	botAuthorId := "<@!" + s.State.User.ID + ">"
	message := strings.ToLower(m.Content)
	if strings.HasPrefix(message, botAuthorId) {
		requestMessage := strings.Replace(message, botAuthorId, "", 1)
		intentAction, intentParameters, paulResponse := dialogflow.ParseRequest(dialogflowProjectId, requestMessage, dialogflowLanguageCode)

		if strings.HasPrefix(intentAction, "workflow") {
			paulResponse = internal.ExecuteWorkflow(temporalClient, intentAction, intentParameters)
		}
		log.Println("Pauls response: ", paulResponse)
		sendMsg, sendErr := s.ChannelMessageSend(m.ChannelID, paulResponse)
		if sendErr != nil {
			log.Fatalln(sendErr)
		}
		log.Println(sendMsg)
	}
}
