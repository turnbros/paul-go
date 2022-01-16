package cluster_event

import (
	"go.temporal.io/sdk/workflow"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
	"log"
	"paul-go/internal/workflows/cluster-event/activities"
)

func ClusterEventMessage(ctx workflow.Context, eventOp watch.EventType, event v1.Event) error {

	var activityErr error
	switch eventOp {
	case watch.Added:
		execErr := workflow.ExecuteActivity(ctx, activities.AddEventMessage, event).Get(ctx, &activityErr)
		if execErr != nil {
			log.Fatalln("ADD activity execution failed: ", execErr)
		}
	case watch.Modified:
		execErr := workflow.ExecuteActivity(ctx, activities.UpdateEventMessage, event).Get(ctx, &activityErr)
		if execErr != nil {
			log.Fatalln("Modify activity execution failed: ", execErr)
		}
	case watch.Deleted:
		execErr := workflow.ExecuteActivity(ctx, activities.RemoveEventMessage, event).Get(ctx, &activityErr)
		if execErr != nil {
			log.Fatalln("Delete activity execution failed: ", execErr)
		}
	}

	return activityErr
}
