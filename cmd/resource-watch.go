package main

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"paul/internal"
	"paul/internal/util"
	dialogflow_entity "paul/internal/workflows/dialogflow-entity"
	dialogflow_entity_util "paul/internal/workflows/dialogflow-entity/util"
	"time"
)

var namespaceEntityTypeId = "47bbdbd2-75a4-42c0-a091-7aaf0aae12e9"
var serviceEntityTypeId = "8f548cae-619c-45b9-8c88-813d68e75135"
var deploymentEntityTypeId = "450b06e5-3fa0-41c4-914c-dfa2bece48d0"
var podEntityTypeId = "7ed95939-23ff-4ab9-bc9e-7f2a0dcc23d6"
var nodeEntityTypeId = "1cdfdd7e-e6b9-422b-bd65-9612157e7500"

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

	// Get a list of namespaces
	namespaceList, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	var namespaceNames []string
	for _, namespace := range namespaceList.Items {
		namespaceNames = append(namespaceNames, namespace.Name)
	}
	setEntity(temporalClient, namespaceEntityTypeId, namespaceNames)

	// Get a list of services
	serviceList, err := client.CoreV1().Services(v1.NamespaceAll).List(ctx, metav1.ListOptions{})
	var serviceNames []string
	for _, service := range serviceList.Items {
		serviceNames = append(serviceNames, service.Name)
	}
	setEntity(temporalClient, serviceEntityTypeId, serviceNames)

	// Get a list of Deployment
	deploymentApi := client.AppsV1().Deployments(v1.NamespaceAll)
	deploymentList, err := deploymentApi.List(ctx, metav1.ListOptions{})
	var deploymentNames []string
	for _, deployment := range deploymentList.Items {
		deploymentNames = append(deploymentNames, deployment.Name)
	}
	setEntity(temporalClient, deploymentEntityTypeId, deploymentNames)

	// Get a list of pods
	podList, err := client.CoreV1().Pods(v1.NamespaceAll).List(ctx, metav1.ListOptions{})
	var podNames []string
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}
	setEntity(temporalClient, podEntityTypeId, podNames)

	namespaceWatcher, err := client.CoreV1().Namespaces().Watch(ctx, metav1.ListOptions{ResourceVersion: resourceVersion})
	serviceWatcher, err := client.CoreV1().Services(v1.NamespaceDefault).Watch(ctx, metav1.ListOptions{ResourceVersion: resourceVersion})
	deploymentWatcher, err := client.CoreV1().Services(v1.NamespaceDefault).Watch(ctx, metav1.ListOptions{ResourceVersion: resourceVersion})
	podWatcher, err := client.CoreV1().Pods(v1.NamespaceAll).Watch(ctx, metav1.ListOptions{ResourceVersion: resourceVersion})

	if err != nil {
		log.Fatal(err)
	}

	go watchNamespaces(temporalClient, namespaceWatcher)
	go watchServices(temporalClient, serviceWatcher)
	go watchDeployments(temporalClient, deploymentWatcher)
	go watchPods(temporalClient, podWatcher)

	log.Println("Waiting for events.")
	for {
		time.Sleep(time.Second)
	}
}

func watchNamespaces(temporalClient client.Client, watcher watch.Interface) {
	entityTypeId := namespaceEntityTypeId
	for event := range watcher.ResultChan() {
		ns := event.Object.(*v1.Namespace)
		switch event.Type {
		case watch.Added:
			fmt.Printf("Namespace %s added\n", ns.ObjectMeta.Name)
			addEntity(temporalClient, entityTypeId, []string{ns.ObjectMeta.Name})
		case watch.Deleted:
			fmt.Printf("Namespace %s deleted\n", ns.ObjectMeta.Name)
			addEntity(temporalClient, entityTypeId, []string{ns.ObjectMeta.Name})
		}
	}
}

func watchServices(temporalClient client.Client, watcher watch.Interface) {
	entityTypeId := serviceEntityTypeId
	for event := range watcher.ResultChan() {
		svc := event.Object.(*v1.Service)
		switch event.Type {
		case watch.Added:
			fmt.Printf("Service %s/%s added\n", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
			addEntity(temporalClient, entityTypeId, []string{svc.ObjectMeta.Name})
		case watch.Deleted:
			fmt.Printf("Service %s/%s deleted", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
			removeEntity(temporalClient, entityTypeId, []string{svc.ObjectMeta.Name})
		}
	}
}

func watchDeployments(temporalClient client.Client, watcher watch.Interface) {
	entityTypeId := podEntityTypeId
	for event := range watcher.ResultChan() {
		deployment := event.Object.(*appsV1.Deployment)
		switch event.Type {
		case watch.Added:
			fmt.Printf("Deployment %s/%s added\n", deployment.ObjectMeta.Namespace, deployment.ObjectMeta.Name)
			addEntity(temporalClient, entityTypeId, []string{deployment.ObjectMeta.Name})
		case watch.Deleted:
			fmt.Printf("Deployment %s/%s deleted\n", deployment.ObjectMeta.Namespace, deployment.ObjectMeta.Name)
			removeEntity(temporalClient, entityTypeId, []string{deployment.ObjectMeta.Name})
		}
	}
}

func watchPods(temporalClient client.Client, watcher watch.Interface) {
	entityTypeId := podEntityTypeId
	for event := range watcher.ResultChan() {
		pod := event.Object.(*v1.Pod)
		switch event.Type {
		case watch.Added:
			fmt.Printf("Pod %s/%s added\n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
			addEntity(temporalClient, entityTypeId, []string{pod.ObjectMeta.Name})
		case watch.Deleted:
			fmt.Printf("Pod %s/%s deleted\n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
			removeEntity(temporalClient, entityTypeId, []string{pod.ObjectMeta.Name})
		}
	}
}

func setEntity(temporalClient client.Client, typeId string, values []string) {
	updateEntity(temporalClient, dialogflow_entity_util.SET, typeId, values)
}
func addEntity(temporalClient client.Client, typeId string, values []string) {
	updateEntity(temporalClient, dialogflow_entity_util.ADD, typeId, values)
}
func removeEntity(temporalClient client.Client, typeId string, values []string) {
	updateEntity(temporalClient, dialogflow_entity_util.REMOVE, typeId, values)
}
func updateEntity(temporalClient client.Client, opCode dialogflow_entity_util.EntityOP, typeId string, values []string) {
	entityRequest := dialogflow_entity_util.EntityRequest{
		Operation:    opCode,
		EntityType:   typeId,
		EntityValues: values,
	}
	_ = dialogflow_entity.ExecuteWorkflow(temporalClient, entityRequest)
}
