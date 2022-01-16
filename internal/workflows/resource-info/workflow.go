package resource_info

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"paul-go/internal/workflows/resource-info/activities"
	rs "paul-go/internal/workflows/resource-info/util"
	"time"
)

func GetResourceInfo(ctx workflow.Context, resourceRequest rs.ResourceRequest) (*string, error) {

	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    2,
	}
	activityOptions := workflow.ActivityOptions{
		RetryPolicy:         retryPolicy,
		StartToCloseTimeout: 2 * time.Minute,
	}

	var response string
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	if resourceRequest.RequestType == "count" {
		switch resourceRequest.ResourceType {
		case "namespace":
			err := workflow.ExecuteActivity(ctx, activities.CountNamespaces, resourceRequest).Get(ctx, &response)
			if err != nil {
				return nil, err
			}
		case "pod":
			err := workflow.ExecuteActivity(ctx, activities.CountPods, resourceRequest).Get(ctx, &response)
			if err != nil {
				return nil, err
			}
		case "service":
			err := workflow.ExecuteActivity(ctx, activities.CountServices, resourceRequest).Get(ctx, &response)
			if err != nil {
				return nil, err
			}
		}
	}

	if resourceRequest.RequestType == "list" {
		switch resourceRequest.ResourceType {
		case "pod":
			err := workflow.ExecuteActivity(ctx, activities.ListPods, resourceRequest).Get(ctx, &response)
			if err != nil {
				return nil, err
			}
		}
	}

	return &response, nil
}
