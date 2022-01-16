package dialogflow_entity

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"log"
	entityActivities "paul-go/internal/workflows/dialogflow-entity/activities"
	util "paul-go/internal/workflows/dialogflow-entity/util"
	"time"
)

func UpdateEntityType(ctx workflow.Context, request util.EntityRequest) error {
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

	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	var response error
	switch request.Operation {
	case util.SET:
		err := workflow.ExecuteActivity(ctx, entityActivities.SetEntityValue, request).Get(ctx, &response)
		if err != nil {
			log.Fatalln("SET activity execution failed: ", err)
		}
	case util.ADD:
		err := workflow.ExecuteActivity(ctx, entityActivities.AddEntityValue, request).Get(ctx, &response)
		if err != nil {
			log.Fatalln("ADD activity execution failed: ", err)
		}
	case util.REMOVE:
		err := workflow.ExecuteActivity(ctx, entityActivities.RemoveEntityValue, request).Get(ctx, &response)
		if err != nil {
			log.Fatalln("REMOVE activity execution failed: ", err)
		}
	}
	return response
}
