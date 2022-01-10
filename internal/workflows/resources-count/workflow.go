package resources_count

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"log"
	"time"
)

const TaskQueue = "ResourceCount"

func StartWorker(client client.Client) {
	workerOptions := worker.Options{}
	workerBee := worker.New(client, TaskQueue, workerOptions)

	workerBee.RegisterWorkflow(CountResources)
	workerBee.RegisterActivity(CountAll)
	// workerBee.RegisterActivity(CountPods)
	// workerBee.RegisterActivity(CountServices)
	// workerBee.RegisterActivity(CountIngresses)

	err := workerBee.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func ExecuteWorkflow(clientSession client.Client, requestParameters string) client.WorkflowRun {

	// Setup the workflow options.
	// TODO: Maybe we could store workflow execution settings in configmap
	workflowOptions := client.StartWorkflowOptions{
		ID:        "resources-count_" + uuid.New().String(),
		TaskQueue: TaskQueue,
	}

	// Unmarshall the dialogflow queryResult parameters into a CountRequest object
	countRequest := CountRequest{} //make(map[string]CountRequest)
	err := json.Unmarshal([]byte(requestParameters), &countRequest)
	if err != nil {
		log.Fatalln("Failed to marshall the request parameters")
		panic(err)
	}

	// kick off the workflow and
	workExec, err := clientSession.ExecuteWorkflow(context.Background(), workflowOptions, CountResources, countRequest)
	if err != nil {
		log.Fatalln("Failed to execute workflow: ", err)
		panic(err)
	}
	return workExec
}

func CountResources(ctx workflow.Context, countRequest CountRequest) (*string, error) {
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
	var asdf string
	err := workflow.ExecuteActivity(ctx, CountAll).Get(ctx, &asdf)
	if err != nil {
		return &asdf, err
	}

	fmt.Println("the result was ", asdf)

	return &asdf, nil
}
