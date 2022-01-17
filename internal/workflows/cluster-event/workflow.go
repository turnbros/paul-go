package cluster_event

import (
	"go.temporal.io/sdk/workflow"
	"log"
	"paul-go/internal/util"
	"paul-go/internal/workflows/cluster-event/activities"
	"time"
)

const eventChannelID = "932458855999352902"

func ClusterEventMessage(ctx workflow.Context, event *util.ClusterEventMessage) error {
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
	execErr := workflow.ExecuteActivity(ctx, activities.AddEventMessage, eventChannelID, event).Get(ctx, &messageId)
	if execErr != nil {
		log.Fatalln("ADD activity execution failed: ", execErr)
	}

	objId2MsgId[event.EventUID] = messageId

	log.Println("Getting ready to wait for signal...")
	var signalName = "EVENT_MODIFIED"
	var eventExists = true
	for eventExists {
		signalChan := workflow.GetSignalChannel(ctx, signalName)
		signal := workflow.NewSelector(ctx)
		signal.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
			log.Println("Add the receive function")
			var signalVal util.ClusterEventMessage
			c.Receive(ctx, &signalVal)
			workflow.GetLogger(ctx).Info("Received signal!", "Signal", signalName, "value", signalVal.EventMessage)
			log.Println("Assign the object uid")
			eventUID := &signalVal.EventUID

			switch event.EventType {
			case "MODIFIED":
				log.Println("Event was modified")
				activityErr = eventModified(ctx, objId2MsgId[*eventUID], *event)
			case "DELETED":
				log.Println("Event was deleted")
				activityErr = eventDeleted(ctx, objId2MsgId[*eventUID])
				eventExists = false
			}
		})
		log.Println("the selector")
		signal.Select(ctx)
	}

	log.Println("All done with workflow")
	return activityErr
}

func eventModified(ctx workflow.Context, discordMsgID string, event util.ClusterEventMessage) error {
	var activityErr error
	execErr := workflow.ExecuteActivity(ctx, activities.UpdateEventMessage, eventChannelID, discordMsgID, event).Get(ctx, &activityErr)
	if execErr != nil {
		log.Fatalln("Modify activity execution failed: ", execErr)
	}
	return activityErr
}

func eventDeleted(ctx workflow.Context, discordMsgID string) error {
	var activityErr error
	execErr := workflow.ExecuteActivity(ctx, activities.RemoveEventMessage, eventChannelID, discordMsgID).Get(ctx, &activityErr)
	if execErr != nil {
		log.Fatalln("Delete activity execution failed: ", execErr)
	}
	return activityErr
}
