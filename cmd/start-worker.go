package main

import (
	"flag"
	"log"
	"os"
	"paul-go/internal"
	dialogflow_entity "paul-go/internal/workflows/dialogflow-entity"
	resourceInfo "paul-go/internal/workflows/resource-info"
)

func main() {

	workflow := flag.String("workflow", "", "The workflow this worker will start")
	flag.Parse()
	if *workflow == "" {
		log.Println("Failed to start worker: -workflow missing. Please run start-worker.go --help for more information")
		os.Exit(2)
	}

	log.Println("Connecting to Temporal...")
	temporalClient := internal.StartTemporal()
	defer temporalClient.Close()

	log.Println("Trying to start worker for workflow: ", *workflow)
	switch *workflow {

	case "resource-info":
		resourceInfo.StartWorker(temporalClient)

	case "dialogflow-entity":
		dialogflow_entity.StartWorker(temporalClient)

	default:
		log.Fatalln("What the hell is ", *workflow, "? I've never heard of that workflow!")
		os.Exit(3)
	}
}
