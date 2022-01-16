package cluster_event

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
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

func ExecuteWorkflow(clientSession client.Client, eventOp watch.EventType, event v1.Event) {
	workflowID := fmt.Sprintf("cluster-event-%v", event.UID)

	if eventOp != watch.Added {
		err := clientSession.SignalWorkflow(context.Background(), workflowID, "", string(eventOp), event)
		if err != nil {
			log.Fatalln("Error signaling client", err)
		}
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: TaskQueue,
	}
	_, err := clientSession.ExecuteWorkflow(context.Background(), workflowOptions, ClusterEventMessage, eventOp, event)
	if err != nil {
		log.Fatalln("Failed to execute workflow: ", err)
	}
}
