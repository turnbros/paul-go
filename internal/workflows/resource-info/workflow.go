package resource_info

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"log"
	"time"
)

const TaskQueue = "ResourceInfo"

func StartWorker(client client.Client) {
	workerOptions := worker.Options{}
	workerBee := worker.New(client, TaskQueue, workerOptions)
	workerBee.RegisterWorkflow(GetResourceInfo)
	workerBee.RegisterActivity(ListPods)

	err := workerBee.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func GetResourceInfo(ctx workflow.Context, infoRequest InfoRequest) (*string, error) {
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    2,
	}
	activityOptions := workflow.ActivityOptions{
		RetryPolicy:         retryPolicy,
		StartToCloseTimeout: 2 * time.Minute,
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var response string
	switch infoRequest.ResourceType {
	case "pod":
		err := workflow.ExecuteActivity(ctx, ListPods, infoRequest).Get(ctx, &response)
		if err != nil {
			return nil, err
		}
	}

	return &response, nil
}

func ExecuteWorkflow(clientSession client.Client, requestParameters string) client.WorkflowRun {

	// Setup the workflow options.
	// TODO: Maybe we could store workflow execution settings in configmap
	workflowOptions := client.StartWorkflowOptions{
		ID:        "resources-info_" + uuid.New().String(),
		TaskQueue: TaskQueue,
	}

	// Unmarshall the dialogflow queryResult parameters into a CountRequest object
	infoRequest := InfoRequest{} //make(map[string]CountRequest)
	err := json.Unmarshal([]byte(requestParameters), &infoRequest)
	if err != nil {
		log.Fatalln("Failed to marshall the request parameters")
		panic(err)
	}

	// kick off the workflow and
	workExec, err := clientSession.ExecuteWorkflow(context.Background(), workflowOptions, GetResourceInfo, infoRequest)
	if err != nil {
		log.Fatalln("Failed to execute workflow: ", err)
		panic(err)
	}
	return workExec
}
