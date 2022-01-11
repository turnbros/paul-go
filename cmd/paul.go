package main

import (
	"encoding/json"
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

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Create a message log entry
	messageLogEntry := MessageLogEntry{
		Timestamp:      string(m.Timestamp),
		MsgID:          m.ID,
		MsgChannelID:   m.ChannelID,
		MsgGuildID:     m.GuildID,
		MsgContent:     m.Content,
		AuthorID:       m.Author.ID,
		AuthorBot:      m.Author.Bot,
		AuthorEmail:    m.Author.Email,
		AuthorUsername: m.Author.Username,
	}

	log.Println(json.Marshal(messageLogEntry))

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
