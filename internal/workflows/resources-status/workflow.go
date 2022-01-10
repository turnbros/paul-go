package resources_status

import (
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
