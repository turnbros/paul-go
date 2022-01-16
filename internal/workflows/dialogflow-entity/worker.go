package dialogflow_entity

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"paul-go/internal/workflows/dialogflow-entity/activities"
	"paul-go/internal/workflows/dialogflow-entity/util"
)

const TaskQueue = "DialogflowEntity"

func StartWorker(client client.Client) {
	workerOptions := worker.Options{
		WorkerActivitiesPerSecond: 5,
	}
	workerBee := worker.New(client, TaskQueue, workerOptions)
	workerBee.RegisterWorkflow(UpdateEntityType)

	workerBee.RegisterActivity(activities.SetEntityValue)
	workerBee.RegisterActivity(activities.AddEntityValue)
	workerBee.RegisterActivity(activities.RemoveEntityValue)

	err := workerBee.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func ExecuteWorkflow(clientSession client.Client, entityRequest util.EntityRequest) client.WorkflowRun {

	// Setup the workflow options.
	// TODO: Maybe we could store workflow execution settings in configmap
	workflowOptions := client.StartWorkflowOptions{
		ID:        fmt.Sprintf("dialogflow-entity-%v_%v", entityRequest.Operation, uuid.New().String()),
		TaskQueue: TaskQueue,
	}

	// kick off the workflow
	workExec, err := clientSession.ExecuteWorkflow(context.Background(), workflowOptions, UpdateEntityType, entityRequest)
	if err != nil {
		log.Fatalln("Failed to execute workflow: ", err)
		panic(err)
	}
	return workExec
}
