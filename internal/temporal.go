package internal

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"log"
	"os"
	"paul-go/internal/util"
	resource_info "paul-go/internal/workflows/resource-info"
)

func StartTemporal() client.Client {
	temporalConfig := util.GetTemporalConfig()

	temporalHostPort := os.Getenv("TEMPORAL_HOSTPORT")
	if temporalHostPort == "" {
		temporalHostPort = fmt.Sprintf("%v:%v", temporalConfig["host"], temporalConfig["port"])
	}

	temporalNamespace := os.Getenv("TEMPORAL_NAMESPACE")
	if temporalNamespace == "" {
		temporalNamespace = fmt.Sprintf("%v", temporalConfig["namespace"])
	}

	log.Println(fmt.Sprintf("%v:%v", temporalConfig["host"], temporalConfig["port"]))
	temporalClient, err := client.NewClient(client.Options{
		// HostPort:  fmt.Sprintf("%v:%v", temporalConfig["host"], temporalConfig["port"]),
		// Namespace: fmt.Sprintf("%v", temporalConfig["namespace"]),
		HostPort:  temporalHostPort,
		Namespace: temporalNamespace,
	})
	if err != nil {
		panic(err)
	}

	return temporalClient
}

func ExecuteWorkflow(temporalClient client.Client, intentAction string, intentParameters string) string {

	var workExec client.WorkflowRun
	var executionResponse string

	switch intentAction {
	case "workflow.resource_info":
		workExec = resource_info.ExecuteWorkflow(temporalClient, intentParameters)
	case "workflow.resource_count":
		workExec = resource_info.ExecuteWorkflow(temporalClient, intentParameters)
	case "workflow.resource_status":
		workExec = resource_info.ExecuteWorkflow(temporalClient, intentParameters)
	default:
		panic("can't find workflow: " + intentAction)
	}

	log.Println("Started workflow", "WorkflowID", workExec.GetID(), "RunID", workExec.GetRunID())
	err := workExec.Get(context.Background(), &executionResponse)
	if err != nil {
		panic(err)
	}
	return executionResponse
}
