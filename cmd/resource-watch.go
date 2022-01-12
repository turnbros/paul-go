package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"paul/internal"
	"paul/internal/util"
	dialogflow_entity "paul/internal/workflows/dialogflow-entity"
	"time"
)

func main() {

	log.Println("Connecting to Temporal...")
	temporalClient := internal.StartTemporal()
	defer temporalClient.Close()

	client := util.GetKubeClient()
	ctx := context.Background()

	var api = client.CoreV1().Endpoints(v1.NamespaceAll)
	endpoints, err := api.List(ctx, metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	resourceVersion := endpoints.ListMeta.ResourceVersion

	namespaceWatcher, err := client.CoreV1().Namespaces().Watch(ctx, metav1.ListOptions{ResourceVersion: resourceVersion})
	serviceWatcher, err := client.CoreV1().Services(v1.NamespaceDefault).Watch(ctx, metav1.ListOptions{ResourceVersion: resourceVersion})
	podWatcher, err := client.CoreV1().Pods(v1.NamespaceAll).Watch(ctx, metav1.ListOptions{ResourceVersion: resourceVersion})

	if err != nil {
		log.Fatal(err)
	}

	go watchNamespaces(namespaceWatcher)
	go watchServices(serviceWatcher)
	go watchPods(podWatcher)

	log.Println("Waiting for events.")
	for {
		time.Sleep(time.Second)
	}
}

func watchNamespaces(watcher watch.Interface) {
	for event := range watcher.ResultChan() {
		svc := event.Object.(*v1.Namespace)
		switch event.Type {
		case watch.Added:
			fmt.Printf("Namespace %s added\n", svc.ObjectMeta.Name)
			//		case watch.Modified:
			//			fmt.Printf("Service %s/%s modified", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
			//		case watch.Deleted:
			//			fmt.Printf("Service %s/%s deleted", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
		}
	}
}

func watchServices(watcher watch.Interface) {
	for event := range watcher.ResultChan() {
		svc := event.Object.(*v1.Service)
		switch event.Type {
		case watch.Added:
			fmt.Printf("Service %s/%s added\n", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
			//		case watch.Modified:
			//			fmt.Printf("Service %s/%s modified", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
			//		case watch.Deleted:
			//			fmt.Printf("Service %s/%s deleted", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
		}
	}
}

func watchPods(watcher watch.Interface) {
	for event := range watcher.ResultChan() {
		pod := event.Object.(*v1.Pod)
		switch event.Type {
		case watch.Added:
			fmt.Printf("Pod %s/%s added\n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
			entityRequest := dialogflow_entity.EntityRequest{
				Operation:    dialogflow_entity.ADD,
				EntityType:   "05519378-3426-4598-8d91-4a01fbb0d2a8",
				EntityValues: []string{pod.ObjectMeta.Name},
			}
			_ = dialogflow_entity.ExecuteWorkflow(temporalClient, entityRequest)
		case watch.Deleted:
			fmt.Printf("Pod %s/%s deleted\n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
			entityRequest := dialogflow_entity.EntityRequest{
				Operation:    dialogflow_entity.REMOVE,
				EntityType:   "05519378-3426-4598-8d91-4a01fbb0d2a8",
				EntityValues: []string{pod.ObjectMeta.Name},
			}
			_ = dialogflow_entity.ExecuteWorkflow(temporalClient, entityRequest)
		}
	}
}
