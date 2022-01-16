package resource_info

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"paul-go/internal/workflows/resource-info/activities"
	rs "paul-go/internal/workflows/resource-info/util"
)

const TaskQueue = "ResourceInfo"

func StartWorker(client client.Client) {
	workerOptions := worker.Options{}
	workerBee := worker.New(client, TaskQueue, workerOptions)
	workerBee.RegisterWorkflow(GetResourceInfo)

	workerBee.RegisterActivity(activities.CountPods)
	workerBee.RegisterActivity(activities.CountServices)
	workerBee.RegisterActivity(activities.CountNamespaces)
	workerBee.RegisterActivity(activities.ListPods)

	err := workerBee.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func ExecuteWorkflow(clientSession client.Client, requestParameters string) client.WorkflowRun {

	// Setup the workflow options.
	// TODO: Maybe we could store workflow execution settings in configmap
	workflowOptions := client.StartWorkflowOptions{
		ID:        "resources-info_" + uuid.New().String(),
		TaskQueue: TaskQueue,
	}

	// Unmarshall the dialogflow queryResult parameters into a CountRequest object
	resourceRequest := rs.ResourceRequest{} //make(map[string]CountRequest)
	err := json.Unmarshal([]byte(requestParameters), &resourceRequest)
	if err != nil {
		log.Fatalln("Failed to marshall the request parameters")
		panic(err)
	}

	// kick off the workflow and
	workExec, err := clientSession.ExecuteWorkflow(context.Background(), workflowOptions, GetResourceInfo, resourceRequest)
	if err != nil {
		log.Fatalln("Failed to execute workflow: ", err)
		panic(err)
	}
	return workExec
}
