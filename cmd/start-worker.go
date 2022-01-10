package main

import (
	"flag"
	"log"
	"os"
	"paul/internal"
	resourceCount "paul/internal/workflows/resource-count"
	resourceInfo "paul/internal/workflows/resource-info"
	resourceStatus "paul/internal/workflows/resource-status"
)

func main() {

	workflow := flag.String("workflow", "", "The workflow this worker will start")
	flag.Parse()
	if *workflow == "" {
		log.Println("Failed to start worker: -workflow missing. Please run start-worker.go --help for more information")
		os.Exit(2)
	}

	log.Println("Connecting to Temporal...")
	temporalClient = internal.StartTemporal()
	defer temporalClient.Close()

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
