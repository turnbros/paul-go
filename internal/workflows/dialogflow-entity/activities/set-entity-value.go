package activities

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"fmt"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"log"
	dialogflow_entity "paul/internal/workflows/dialogflow-entity/util"
)

func SetEntityValue(ctx context.Context, request dialogflow_entity.EntityRequest) error {
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

	var entities []*dialogflowpb.EntityType_Entity
	for _, entityValue := range request.EntityValues {
		entities = append(entities, &dialogflowpb.EntityType_Entity{
			Value: entityValue,
		})
	}

	resourceTypeEntity.Entities = entities

	_, updateError := entityClient.UpdateEntityType(ctx, &dialogflowpb.UpdateEntityTypeRequest{EntityType: resourceTypeEntity, LanguageCode: "en"})
	if updateError != nil {
		log.Fatalln(fmt.Sprintf("Failed to set the values for entity %v: %v ", resourceTypeEntity.Name, updateError))
	} else {
		log.Println(fmt.Sprintf("Value %v has been successfully set for Entity %v", request.EntityValues, resourceTypeEntity.Name))
	}

	return nil
}
