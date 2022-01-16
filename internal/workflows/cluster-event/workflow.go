package cluster_event

import (
	"go.temporal.io/sdk/workflow"
	"log"
	"paul-go/internal/util"
	"paul-go/internal/workflows/cluster-event/activities"
	"time"
)

const importantEventsChannelID = "931301531179966515"
const normalEventsChannelID = "931301737028001802"
const testEventsChannelID = "932115780768759878"

func ClusterEventMessage(ctx workflow.Context, eventOp string, event util.ClusterEventMessage) error {
	log.Println("Starting ClusterEventMessage...")

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Second * 60,
		ScheduleToStartTimeout: time.Second * 60,
		StartToCloseTimeout:    time.Second * 60,
		HeartbeatTimeout:       time.Second * 10,
		WaitForCancellation:    false,
	})

	// This ain't right. Needs to create or wait for signal.
	objId2MsgId := make(map[string]string)

	// Create the Discord message
	var activityErr error
	var messageId string
	execErr := workflow.ExecuteActivity(ctx, activities.AddEventMessage, testEventsChannelID, event).Get(ctx, &messageId)
	if execErr != nil {
		log.Fatalln("ADD activity execution failed: ", execErr)
	}
	objId2MsgId[event.ObjectUID] = messageId

	/*	switch eventOp {
		case "ADDED": // watch.Added:
			var messageId string
			execErr := workflow.ExecuteActivity(ctx, activities.AddEventMessage, testEventsChannelID, event).Get(ctx, &messageId)
			if execErr != nil {
				log.Fatalln("ADD activity execution failed: ", execErr)
			}
			objId2MsgId[event.ObjectUID] = messageId

		case "MODIFIED":
			execErr := workflow.ExecuteActivity(ctx, activities.UpdateEventMessage, testEventsChannelID, objId2MsgId[event.ObjectUID], event).Get(ctx, &activityErr)
			if execErr != nil {
				log.Fatalln("Modify activity execution failed: ", execErr)
			}

		case "DELETED":
			execErr := workflow.ExecuteActivity(ctx, activities.RemoveEventMessage, testEventsChannelID, objId2MsgId[event.ObjectUID]).Get(ctx, &activityErr)
			if execErr != nil {
				log.Fatalln("Delete activity execution failed: ", execErr)
			}

		}*/

	var signalName = "MODIFIED"
	modifiedSignalChan := workflow.GetSignalChannel(ctx, signalName)
	sModified := workflow.NewSelector(ctx)
	sModified.AddReceive(modifiedSignalChan, func(c workflow.ReceiveChannel, more bool) {
		var signalVal util.ClusterEventMessage
		c.Receive(ctx, &signalVal)
		workflow.GetLogger(ctx).Info("Received signal!", "Signal", signalName, "value", signalVal.EventMessage)
		execErr := workflow.ExecuteActivity(ctx, activities.UpdateEventMessage, testEventsChannelID, objId2MsgId[signalVal.ObjectUID], signalVal).Get(ctx, &activityErr)
		if execErr != nil {
			log.Fatalln("Modify activity execution failed: ", execErr)
		}
	})

	for {
		sModified.Select(ctx)
	}

	log.Println("All done with workflow")
	return activityErr
}

func eventMessageModified(ctx workflow.Context, eventOp string, event util.ClusterEventMessage) error {
	return nil
}
