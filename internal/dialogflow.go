package internal

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"log"
)

func ParseRequest(projectID, userRequest, languageCode string) (string, string, string) {
	detectedIntent, err := detectIntentText(projectID, uuid.New().String(), userRequest, languageCode)

	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	intentAction := detectedIntent.GetAction()
	log.Println("Intent Action: ", intentAction)
	intentParameters, _ := json.Marshal(detectedIntent.GetParameters().AsMap())
	log.Println("Intent Parameters: ", string(intentParameters))

	return intentAction, string(intentParameters), detectedIntent.FulfillmentText
}

func detectIntentText(projectID, sessionID, text, languageCode string) (*dialogflowpb.QueryResult, error) {
	ctx := context.Background()

	sessionClient, err := dialogflow.NewSessionsClient(ctx)
	if err != nil {
		return nil, err
	}
	defer sessionClient.Close()

	if projectID == "" || sessionID == "" {
		return nil, errors.New(fmt.Sprintf("Received empty project (%s) or session (%s)", projectID, sessionID))
	}

	sessionPath := fmt.Sprintf("projects/%s/agent/sessions/%s", projectID, sessionID)
	textInput := dialogflowpb.TextInput{Text: text, LanguageCode: languageCode}
	queryTextInput := dialogflowpb.QueryInput_Text{Text: &textInput}
	queryInput := dialogflowpb.QueryInput{Input: &queryTextInput}
	request := dialogflowpb.DetectIntentRequest{Session: sessionPath, QueryInput: &queryInput}

	response, err := sessionClient.DetectIntent(ctx, &request)
	if err != nil {
		return nil, err
	}

	return response.GetQueryResult(), nil
}
