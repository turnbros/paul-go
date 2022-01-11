package resource_count

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"log"
	"regexp"
	"strings"
	"time"
)

const TaskQueue = "ResourceCount"

func StartWorker(client client.Client) {
	workerOptions := worker.Options{}
	workerBee := worker.New(client, TaskQueue, workerOptions)

	workerBee.RegisterWorkflow(GetResourceCount)
	workerBee.RegisterActivity(CountAll)
	workerBee.RegisterActivity(CountPods)
	workerBee.RegisterActivity(CountServices)
	workerBee.RegisterActivity(CountNamespaces)

	err := workerBee.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func GetResourceCount(ctx workflow.Context, countRequest CountRequest) (*string, error) {
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

	if countRequest.ResourceScope != "" {
		log.Println(countRequest.ResourceScope)
		// Remove the spaces that dialog flow seems hellbent on adding.
		countRequest.ResourceScope = strings.Replace(countRequest.ResourceScope, " ", "", -1)

		// Lean up the resource type to make sure whatever we get is in a valid namespace format
		re := regexp.MustCompile("^[a-zA-Z0-9\\-]{1,63}")
		match := re.FindStringSubmatch(countRequest.ResourceScope)
		log.Println(match)
		if len(match) > 1 {
			countRequest.ResourceScope = match[0]
		}
	}

	switch countRequest.ResourceType {
	case "namespace":
		err := workflow.ExecuteActivity(ctx, CountNamespaces, countRequest).Get(ctx, &response)
		if err != nil {
			return nil, err
		}
	case "pod":
		err := workflow.ExecuteActivity(ctx, CountPods, countRequest).Get(ctx, &response)
		if err != nil {
			return nil, err
		}
	case "service":
		err := workflow.ExecuteActivity(ctx, CountServices, countRequest).Get(ctx, &response)
		if err != nil {
			return nil, err
		}
	default:
		err := workflow.ExecuteActivity(ctx, CountAll, countRequest).Get(ctx, &response)
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
	workExec, err := clientSession.ExecuteWorkflow(context.Background(), workflowOptions, GetResourceCount, countRequest)
	if err != nil {
		log.Fatalln("Failed to execute workflow: ", err)
		panic(err)
	}
	return workExec
}
