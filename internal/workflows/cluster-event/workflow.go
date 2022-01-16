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

	var messageId string
	var activityErr error
	objId2MsgId := make(map[string]string)

	// Create the Discord message
	execErr := workflow.ExecuteActivity(ctx, activities.AddEventMessage, testEventsChannelID, event).Get(ctx, &messageId)
	if execErr != nil {
		log.Fatalln("ADD activity execution failed: ", execErr)
	}

	objId2MsgId[event.EventUID] = messageId

	var signalName = "EVENT_MODIFIED"
	for {
		signalChan := workflow.GetSignalChannel(ctx, signalName)
		signal := workflow.NewSelector(ctx)
		signal.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
			var signalVal util.ClusterEventMessage
			c.Receive(ctx, &signalVal)
			workflow.GetLogger(ctx).Info("Received signal!", "Signal", signalName, "value", signalVal.EventMessage)

			switch event.EventType {
			case "MODIFIED":
				activityErr = eventModified(ctx, objId2MsgId[signalVal.ObjectUID], event)
			case "DELETED":
				activityErr = eventDeleted(ctx, objId2MsgId[signalVal.ObjectUID])
			}
		})
		signal.Select(ctx)
	}

	log.Println("All done with workflow")
	return activityErr
}

func eventModified(ctx workflow.Context, discordMsgID string, event util.ClusterEventMessage) error {
	var activityErr error
	execErr := workflow.ExecuteActivity(ctx, activities.UpdateEventMessage, testEventsChannelID, discordMsgID, event).Get(ctx, &activityErr)
	if execErr != nil {
		log.Fatalln("Modify activity execution failed: ", execErr)
	}
	return activityErr
}

func eventDeleted(ctx workflow.Context, discordMsgID string) error {
	var activityErr error
	execErr := workflow.ExecuteActivity(ctx, activities.RemoveEventMessage, testEventsChannelID, discordMsgID).Get(ctx, &activityErr)
	if execErr != nil {
		log.Fatalln("Delete activity execution failed: ", execErr)
	}
	return activityErr
}
