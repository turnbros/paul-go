package main

import (
	"context"
	"flag"
	"go.temporal.io/sdk/client"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"paul-go/internal"
	"paul-go/internal/util"
	"strings"
	"time"
)

func main() {
	eventAction := flag.String("action", "", "action to take")
	eventCount := flag.Int("count", 1, "number of times seen event")
	flag.Parse()

	log.Println("Connecting to Temporal...")
	temporalClient := internal.StartTemporal()
	defer temporalClient.Close()

	log.Println("Action: ", eventAction)
	log.Println("Count: ", eventCount)

	switch strings.ToLower(*eventAction) {
	case "add":
		log.Println("adding event")
		//sendEvent(temporalClient, watch.Added, makeTestEvent(int32(*eventCount)))
	case "update":
		log.Println("editing event")
	//	sendEvent(temporalClient, watch.Modified, makeTestEvent(int32(*eventCount)))
	case "delete":
		log.Println("deleting event")
		//sendEvent(temporalClient, watch.Deleted, makeTestEvent(int32(*eventCount)))

	}

	kubeClient := util.GetKubeClient()
	ctx := context.Background()
	eventWatcher, err := kubeClient.CoreV1().Events(v1.NamespaceAll).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatalln("Failed to watch events: ", err)
	}
	go watchEventsscratch(temporalClient, eventWatcher)

	log.Println("Waiting for events...")
	for {
		time.Sleep(time.Second)
	}
}

func watchEventsscratch(temporalClient client.Client, watcher watch.Interface) {
	log.Println("Is this event?")
	for event := range watcher.ResultChan() {
		clusterEvent := event.Object.(*v1.Event)
		/*clusterEventMessage := util.ClusterEventMessage{
			SourceComponent:       "",
			SourceHost:            "",
			ObjectKind:            "",
			ObjectNamespace:       "",
			ObjectName:            "",
			ObjectUID:             "",
			ObjectAPIVersion:      "",
			ObjectResourceVersion: "",
			ObjectFieldPath:       "",
			EventName:             "",
			EventReason:           "",
			EventMessage:          "",
			EventCount:            0,
			EventType:             "",
			EventFirstTimestamp:   "",
			EventLastTimestamp:    "",
		}*/
		log.Println("Received event: ", clusterEvent.Name)
		//cluster_event.ExecuteWorkflow(temporalClient, event.Type, clusterEventMessage)
	}
}

func sendEvent(temporalClient client.Client, event util.ClusterEventMessage) {
	log.Println("Received event: ", event.EventName)
	//cluster_event.StartWorkflow(temporalClient, "UPDATE", event)
}

func makeTestEvent(count int32) util.ClusterEventMessage {
	return util.ClusterEventMessage{
		SourceComponent:       "kubelet",
		SourceHost:            "dear-dora",
		ObjectKind:            "Pod",
		ObjectNamespace:       "olm",
		ObjectName:            "operatorhubio-catalog-tbs6n",
		ObjectUID:             "dc7fc5a9-eadf-4aca-af48-aed8fa7ee2df",
		ObjectAPIVersion:      "v1",
		ObjectResourceVersion: "108115203",
		ObjectFieldPath:       "spec.containers{registry-server}",
		EventName:             "operatorhubio-catalog-tbs6n.16cad1f241f33f76",
		EventReason:           "Started",
		EventMessage:          "Started container registry-server",
		EventCount:            count,
		EventType:             "Normal",
		EventFirstTimestamp:   "2021-09-16T20:55:03Z",
		EventLastTimestamp:    "2022-01-16T18:41:24Z",
	}
}
