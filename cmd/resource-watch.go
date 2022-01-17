package main

import (
	"context"
	"go.temporal.io/sdk/client"
	appsV1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"paul-go/internal"
	"paul-go/internal/util"
	cluster_event "paul-go/internal/workflows/cluster-event"
	dialogflow_entity "paul-go/internal/workflows/dialogflow-entity"
	dialogflow_entity_util "paul-go/internal/workflows/dialogflow-entity/util"
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

	log.Println("Connecting to Kubernetes...")
	kubeClient := util.GetKubeClient()
	ctx := context.Background()

	// Get a list of namespaces
	namespaceList, err := kubeClient.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	var namespaceNames []string
	for _, namespace := range namespaceList.Items {
		namespaceNames = append(namespaceNames, namespace.Name)
	}
	setEntity(temporalClient, namespaceEntityTypeId, namespaceNames)
	namespaceWatcher, err := kubeClient.CoreV1().Namespaces().Watch(ctx, metav1.ListOptions{ResourceVersion: namespaceList.ListMeta.ResourceVersion})
	if err != nil {
		log.Fatal(err)
	}
	go watchNamespaces(temporalClient, namespaceWatcher)

	// Get a list of services
	serviceList, err := kubeClient.CoreV1().Services(v1.NamespaceAll).List(ctx, metav1.ListOptions{})
	var serviceNames []string
	for _, service := range serviceList.Items {
		serviceNames = append(serviceNames, service.Name)
	}
	setEntity(temporalClient, serviceEntityTypeId, serviceNames)
	serviceWatcher, err := kubeClient.CoreV1().Services(v1.NamespaceAll).Watch(ctx, metav1.ListOptions{ResourceVersion: serviceList.ListMeta.ResourceVersion})
	if err != nil {
		log.Fatal(err)
	}
	go watchServices(temporalClient, serviceWatcher)

	// Get a list of Deployment
	deploymentList, err := kubeClient.AppsV1().Deployments(v1.NamespaceAll).List(ctx, metav1.ListOptions{})
	var deploymentNames []string
	for _, deployment := range deploymentList.Items {
		deploymentNames = append(deploymentNames, deployment.Name)
	}
	setEntity(temporalClient, deploymentEntityTypeId, deploymentNames)
	/*deploymentWatcher, err := kubeClient.AppsV1().Deployments(v1.NamespaceAll).Watch(ctx, metav1.ListOptions{ResourceVersion: deploymentList.ListMeta.ResourceVersion})
	if err != nil {
		log.Fatal(err)
	}*/
	//go watchDeployments(temporalClient, deploymentWatcher)

	// Get a list of pods
	podList, err := kubeClient.CoreV1().Pods(v1.NamespaceAll).List(ctx, metav1.ListOptions{})
	var podNames []string
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}
	setEntity(temporalClient, podEntityTypeId, podNames)
	podWatcher, err := kubeClient.CoreV1().Pods(v1.NamespaceAll).Watch(ctx, metav1.ListOptions{ResourceVersion: podList.ListMeta.ResourceVersion})
	if err != nil {
		log.Fatal(err)
	}
	go watchPods(temporalClient, podWatcher)

	// Get a list of events
	//	eventList, err := kubeClient.CoreV1().Events(v1.NamespaceAll).List(ctx, metav1.ListOptions{})
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	// eventWatcher, err := kubeClient.CoreV1().Events(v1.NamespaceAll).Watch(ctx, metav1.ListOptions{ResourceVersion: eventList.ListMeta.ResourceVersion})
	eventWatcher, err := kubeClient.CoreV1().Events(v1.NamespaceAll).Watch(ctx, metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	go watchEvents(temporalClient, eventWatcher)

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
			log.Printf("Namespace %s added\n", ns.ObjectMeta.Name)
			addEntity(temporalClient, entityTypeId, []string{ns.ObjectMeta.Name})
		case watch.Deleted:
			log.Printf("Namespace %s deleted\n", ns.ObjectMeta.Name)
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
			log.Printf("Service %s/%s added\n", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
			addEntity(temporalClient, entityTypeId, []string{svc.ObjectMeta.Name})
		case watch.Deleted:
			log.Printf("Service %s/%s deleted\n", svc.ObjectMeta.Namespace, svc.ObjectMeta.Name)
			removeEntity(temporalClient, entityTypeId, []string{svc.ObjectMeta.Name})
		}
	}
}

func watchDeployments(temporalClient client.Client, watcher watch.Interface) {
	entityTypeId := deploymentEntityTypeId
	for event := range watcher.ResultChan() {
		deployment := event.Object.(*appsV1.Deployment)
		switch event.Type {
		case watch.Added:
			log.Printf("Deployment %s/%s added\n", deployment.ObjectMeta.Namespace, deployment.ObjectMeta.Name)
			addEntity(temporalClient, entityTypeId, []string{deployment.ObjectMeta.Name})
		case watch.Deleted:
			log.Printf("Deployment %s/%s deleted\n", deployment.ObjectMeta.Namespace, deployment.ObjectMeta.Name)
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
			log.Printf("Pod %s/%s added\n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
			addEntity(temporalClient, entityTypeId, []string{pod.ObjectMeta.Name})
		case watch.Deleted:
			log.Printf("Pod %s/%s deleted\n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name)
			removeEntity(temporalClient, entityTypeId, []string{pod.ObjectMeta.Name})
		}
	}
}

func watchEvents(temporalClient client.Client, watcher watch.Interface) {
	for event := range watcher.ResultChan() {
		clusterEvent := event.Object.(*v1.Event)
		if event.Type == watch.Added {
			cluster_event.StartWorkflow(temporalClient, clusterEvent)
		} else {
			cluster_event.UpdateWorkflow(temporalClient, clusterEvent)
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
