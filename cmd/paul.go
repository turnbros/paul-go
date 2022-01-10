package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"paul/internal"
	resources_count "paul/internal/workflows/resources-count"

	"go.temporal.io/sdk/client"
)

const dialogflowProjectId = "paul-fmma"
const dialogflowLanguageCode = "en"

func main() {
	// The client is a heavyweight object that should be created once per process.
	temporalClient, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporalClient.Close()

	userRequest := "How many services are there"

	detectedIntent, err := internal.DetectIntentText(dialogflowProjectId, uuid.New().String(), userRequest, dialogflowLanguageCode)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}

	intentAction := detectedIntent.GetAction()
	log.Println("Intent Action: ", intentAction)
	intentParameters, _ := json.Marshal(detectedIntent.GetParameters().AsMap())
	log.Println("Intent Parameters: ", string(intentParameters))

	workExec := resources_count.ExecuteWorkflow(temporalClient, string(intentParameters))
	log.Println("Started workflow", "WorkflowID", workExec.GetID(), "RunID", workExec.GetRunID())

	var executionResponse string
	err = workExec.Get(context.Background(), &executionResponse)
	if err != nil {
		log.Fatalln("Unable get workflow response", err)
	}

	log.Println("Workflow Response:", executionResponse)
}
