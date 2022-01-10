package util

import (
	"encoding/json"
	"log"
)

func GetTemporalConfig() map[string]interface{} {
	var temporalMap map[string]interface{}
	configMap := GetConfigMapData("paul", "paul-cm")
	log.Println(configMap["temporal"])
	err := json.Unmarshal([]byte(configMap["temporal"]), &temporalMap)
	if err != nil {
		log.Fatalln("Failed to temporalMap or something", err)
	}
	return temporalMap
}

func GetDialogflowConfig() map[string]interface{} {
	var dialogflowMap map[string]interface{}
	configMap := GetConfigMapData("paul", "paul-cm")
	log.Println(configMap["dialogflow"])
	err := json.Unmarshal([]byte(configMap["dialogflow"]), &dialogflowMap)
	if err != nil {
		log.Fatalln("Failed to dialogflowMap or something", err)
	}
	return dialogflowMap
}

func GetWorkflowConfig() map[string]interface{} {
	var workflowMap map[string]interface{}
	configMap := GetConfigMapData("paul", "paul-cm")
	log.Println(configMap["workflows"])
	err := json.Unmarshal([]byte(configMap["workflows"]), &workflowMap)
	if err != nil {
		log.Fatalln("Failed to dialogflowMap or something", err)
	}
	return workflowMap
}
