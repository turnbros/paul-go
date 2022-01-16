package cluster_event

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"paul-go/internal/util"
	"paul-go/internal/workflows/cluster-event/activities"
)

const TaskQueue = "ClusterEvent"

func StartWorker(client client.Client) {
	workerOptions := worker.Options{
		WorkerActivitiesPerSecond: 5,
	}
	workerBee := worker.New(client, TaskQueue, workerOptions)
	workerBee.RegisterWorkflow(ClusterEventMessage)
	workerBee.RegisterActivity(activities.AddEventMessage)
	workerBee.RegisterActivity(activities.RemoveEventMessage)
	workerBee.RegisterActivity(activities.UpdateEventMessage)

	err := workerBee.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}

func UpdateWorkflow(clientSession client.Client, event util.ClusterEventMessage) {

}

func ExecuteWorkflow(clientSession client.Client, event util.ClusterEventMessage) {
	log.Println("Starting Worker ExecuteWorkflow...")
	workflowID := fmt.Sprintf("cluster-event-%v", event.ObjectUID)

	log.Println("Workflow exists, sending signal")
	err := clientSession.SignalWorkflow(context.Background(), workflowID, "", string(eventOp), event)
	if err != nil {
		log.Fatalln("Error signaling client", err)
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: TaskQueue,
	}
	log.Println("Adding workflow...")
	_, err := clientSession.ExecuteWorkflow(context.Background(), workflowOptions, ClusterEventMessage, string(eventOp), event)
	if err != nil {
		log.Fatalln("Failed to execute workflow: ", err)
	}
}
