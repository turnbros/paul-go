package cluster_event

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	v1 "k8s.io/api/core/v1"
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

func StartWorkflow(clientSession client.Client, eventObject *v1.Event) {
	log.Println("Starting Worker ExecuteWorkflow...")

	event := parseClusterEvent(eventObject)
	workflowOptions := client.StartWorkflowOptions{
		ID:        getWorkflowID(event),
		TaskQueue: TaskQueue,
	}
	log.Println("Adding workflow...")
	_, err := clientSession.ExecuteWorkflow(context.Background(), workflowOptions, ClusterEventMessage, &event)
	if err != nil {
		log.Fatalln("Failed to execute workflow: ", err)
	}
}

func UpdateWorkflow(clientSession client.Client, eventObject *v1.Event) {
	log.Println("Workflow exists, sending signal")

	event := parseClusterEvent(eventObject)
	workflowOptions := client.StartWorkflowOptions{
		ID:        getWorkflowID(event),
		TaskQueue: TaskQueue,
	}
	_, err := clientSession.SignalWithStartWorkflow(context.Background(), getWorkflowID(event), "EVENT_MODIFIED", event, workflowOptions, ClusterEventMessage, &event)
	if err != nil {
		log.Fatalln("Error signaling client", err)
	}
}

func getWorkflowID(event *util.ClusterEventMessage) string {
	return fmt.Sprintf("cluster-event-%v", event.EventUID)
}

func parseClusterEvent(event *v1.Event) *util.ClusterEventMessage {
	return &util.ClusterEventMessage{
		SourceComponent:       event.Source.Component,
		SourceHost:            event.Source.Host,
		ObjectKind:            event.InvolvedObject.Kind,
		ObjectNamespace:       event.InvolvedObject.Namespace,
		ObjectName:            event.InvolvedObject.Name,
		ObjectUID:             string(event.InvolvedObject.UID),
		ObjectAPIVersion:      event.InvolvedObject.APIVersion,
		ObjectResourceVersion: event.InvolvedObject.ResourceVersion,
		ObjectFieldPath:       event.InvolvedObject.FieldPath,
		EventName:             event.Name,
		EventReason:           event.Reason,
		EventMessage:          event.Message,
		EventCount:            event.Count,
		EventType:             event.Type,
		EventUID:              string(event.UID),
		EventFirstTimestamp:   event.FirstTimestamp.String(),
		EventLastTimestamp:    event.LastTimestamp.String(),
	}
}
