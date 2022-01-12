package dialogflow

import (
	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"context"
	"fmt"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
	"log"
)

func RemoveResourceTypeValue(resourceType string, resourceName string) {

}

func AddResourceTypeValue(resourceType string, resourceName string) {

	// Create the entity client that we'll use to update Pauls resourceType entities
	ctx := context.Background()
	entityClient, err := dialogflow.NewEntityTypesClient(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer entityClient.Close()

	// Craft the request to get the existing resourceType entity
	entityRequest := dialogflowpb.GetEntityTypeRequest{
		Name: fmt.Sprintf("projects/%v/locations/global/agent/entityTypes/%v", "paul-fmma", resourceType),
	}

	// Query Dialogflow for the entity we'd like to update.
	resourceTypeEntity, err := entityClient.GetEntityType(ctx, &entityRequest)
	if err != nil {
		log.Fatalln("Failed to get session entity type", err)
	}
	log.Println(resourceTypeEntity)

	/*
		##############################
		#########    Note    #########
		##############################
		In these parts `dialogflow.EntityType` is referred to as resourceTypeEntities because `EntityType` is
		actually an object that contains (pointer to) a list of `dialogflowpb.EntityType_Entity` objects
		that make up the actual values we care about.

		Example:
		EntityType 					->	[]*EntityType_Entity
		resourceTypeEntity	->	[]*resourceTypeValue
	*/

	// This is the new resource instance we're going to try associate with an EntityType
	resourceTypeValue := &dialogflowpb.EntityType_Entity{
		Value: resourceName,
	}

	// Here we update the list of resource names (EntityType_Entity) that we've associated with the resource type (EntityType)
	var valueChanged bool
	resourceTypeEntity.Entities, valueChanged = appendIfMissing(resourceTypeEntity.Entities, resourceTypeValue)

	if valueChanged {
		_, updateError := entityClient.UpdateEntityType(ctx, &dialogflowpb.UpdateEntityTypeRequest{EntityType: resourceTypeEntity, LanguageCode: "en"})
		if updateError != nil {
			log.Fatalln(fmt.Sprintf("Failed to update the entity %v: %v ", resourceTypeEntity.Name, updateError))
		} else {
			log.Println(fmt.Sprintf("Value %v has been successfully added to Entity %v", resourceName, resourceTypeEntity.Name))
		}
	} else {
		log.Println(fmt.Sprintf("Value %v already present in %v, nothing to do....", resourceName, resourceTypeEntity.Name))
	}
}
