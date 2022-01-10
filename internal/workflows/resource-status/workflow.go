package resource_status

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"log"
)

const TaskQueue = "ResourceStatus"

func StartWorker(client client.Client) {
	workerOptions := worker.Options{}
	workerBee := worker.New(client, TaskQueue, workerOptions)
	workerBee.RegisterWorkflow(GetResourceStatus)
	//workerBee.RegisterActivity(CountAll)

	err := workerBee.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func GetResourceStatus(ctx workflow.Context) (*int, error) {
	return nil, nil
}

func ExecuteWorkflow(clientSession client.Client, requestParameters string) client.WorkflowRun {

	// Setup the workflow options.
	// TODO: Maybe we could store workflow execution settings in configmap
	workflowOptions := client.StartWorkflowOptions{
		ID:        "resources-count_" + uuid.New().String(),
		TaskQueue: TaskQueue,
	}

	// Unmarshall the dialogflow queryResult parameters into a CountRequest object
	countRequest := StatusRequest{} //make(map[string]CountRequest)
	err := json.Unmarshal([]byte(requestParameters), &countRequest)
	if err != nil {
		log.Fatalln("Failed to marshall the request parameters")
		panic(err)
	}

	// kick off the workflow and
	workExec, err := clientSession.ExecuteWorkflow(context.Background(), workflowOptions, GetResourceStatus, countRequest)
	if err != nil {
		log.Fatalln("Failed to execute workflow: ", err)
		panic(err)
	}
	return workExec
}
