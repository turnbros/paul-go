package util

type ClusterEventMessage struct {
	SourceComponent       string `json:"source_component"`
	SourceHost            string `json:"source_host"`
	ObjectKind            string `json:"object_kind"`
	ObjectNamespace       string `json:"object_namespace"`
	ObjectName            string `json:"object_name"`
	ObjectUID             string `json:"object_uid"`
	ObjectAPIVersion      string `json:"object_api_version"`
	ObjectResourceVersion string `json:"object_resource_version"`
	ObjectFieldPath       string `json:"object_field_path"`
	EventName             string `json:"event_name"`
	EventReason           string `json:"event_reason"`
	EventMessage          string `json:"event_message"`
	EventCount            int32  `json:"event_count"`
	EventType             string `json:"event_type"`
	EventFirstTimestamp   string `json:"event_first_timestamp"`
	EventLastTimestamp    string `json:"event_last_timestamp"`
}
