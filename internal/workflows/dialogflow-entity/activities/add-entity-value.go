package activities

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"fmt"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"log"
	dialogflow_entity "paul/internal/workflows/dialogflow-entity/util"
)

func AddEntityValue(ctx context.Context, request dialogflow_entity.EntityRequest) error {

	// Create the entity client that we'll use to update Pauls resourceType entities
	entityCtx := context.Background()
	entityClient, err := dialogflow.NewEntityTypesClient(entityCtx)
	if err != nil {
		log.Fatalln(err)
	}
	defer entityClient.Close()

	// Craft the request to get the existing resourceType entity
	entityRequest := dialogflowpb.GetEntityTypeRequest{
		Name: fmt.Sprintf("projects/%v/locations/global/agent/entityTypes/%v", "paul-fmma", request.EntityType),
	}

	// Query Dialogflow for the entity we'd like to update.
	resourceTypeEntity, err := entityClient.GetEntityType(ctx, &entityRequest)
	if err != nil {
		log.Fatalln("Failed to get session entity type", err)
	}
	log.Println(resourceTypeEntity)

	valueChanged := false
	for _, entityValue := range request.EntityValues {
		var entityAppended bool
		// Here we update the list of resource names (EntityType_Entity) that we've associated with the resource type (EntityType)
		resourceTypeEntity.Entities, entityAppended = appendIfMissing(resourceTypeEntity.Entities, entityValue)
		if entityAppended {
			valueChanged = true
		}
	}

	if valueChanged {
		_, updateError := entityClient.UpdateEntityType(ctx, &dialogflowpb.UpdateEntityTypeRequest{EntityType: resourceTypeEntity, LanguageCode: "en"})
		if updateError != nil {
			log.Fatalln(fmt.Sprintf("Failed to update the entity %v: %v ", resourceTypeEntity.Name, updateError))
		} else {
			log.Println(fmt.Sprintf("Value %v has been successfully added to Entity %v", request.EntityValues, resourceTypeEntity.Name))
		}
	} else {
		log.Println(fmt.Sprintf("Value %v already present in %v, nothing to do....", request.EntityValues, resourceTypeEntity.Name))
	}

	return nil
}

func appendIfMissing(entitySlice []*dialogflowpb.EntityType_Entity, entityValue string) ([]*dialogflowpb.EntityType_Entity, bool) {
	// To do this we must check all the entities in entitySlice
	for _, entity := range entitySlice {
		if entity.Value == entityValue {
			return entitySlice, false
		}
	}
	entitySlice = append(entitySlice, &dialogflowpb.EntityType_Entity{
		Value: entityValue,
	})
	return entitySlice, true
}
