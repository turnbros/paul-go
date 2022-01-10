package main

import (
	"flag"
	"go.temporal.io/sdk/client"
	"log"
	"os"
	resourceInfo "paul/internal/workflows/resource-info"
	resourceCount "paul/internal/workflows/resources-count"
	resourceStatus "paul/internal/workflows/resources-status"
)

func main() {

	workflow := flag.String("workflow", "", "The workflow this worker will start")
	flag.Parse()
	if *workflow == "" {
		log.Println("Failed to start worker: -workflow missing. Please run start-worker.go --help for more information")
		os.Exit(2)
	}

	log.Println("Connecting to Temporal...")
	clientOptions := client.Options{
		HostPort:  "127.0.0.1:7233",
		Namespace: "default",
	}
	temporalClient, err := client.NewClient(clientOptions)
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
		os.Exit(4)
	} else {
		defer temporalClient.Close()
		log.Println("Connected!")
	}

	log.Println("Trying to start worker for workflow: ", *workflow)
	switch *workflow {

	case "resource-info":
		resourceInfo.StartWorker(temporalClient)

	case "resource-count":
		resourceCount.StartWorker(temporalClient)

	case "resource-status":
		resourceStatus.StartWorker(temporalClient)

	default:
		log.Fatalln("What the hell is ", *workflow, "? I've never heard of that workflow!")
		os.Exit(3)
	}
}
