package internal

import (
	"context"
	"go.temporal.io/sdk/client"
	"log"
	resource_count "paul/internal/workflows/resource-count"
	resource_info "paul/internal/workflows/resource-info"
	resource_status "paul/internal/workflows/resource-status"
)

func StartTemporal() client.Client {
	temporalClient, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
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
	case "resource_info":
		workExec = resource_info.ExecuteWorkflow(temporalClient, intentParameters)
	case "resource_count":
		workExec = resource_count.ExecuteWorkflow(temporalClient, intentParameters)
	case "resource_status":
		workExec = resource_status.ExecuteWorkflow(temporalClient, intentParameters)
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
