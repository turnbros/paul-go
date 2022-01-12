package dialogflow

import dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"

func appendIfMissing(entitySlice []*dialogflowpb.EntityType_Entity, entity *dialogflowpb.EntityType_Entity) ([]*dialogflowpb.EntityType_Entity, bool) {
	for _, ele := range entitySlice {
		if ele.Value == entity.Value {
			return entitySlice, false
		}
	}
	return append(entitySlice, entity), true
}
