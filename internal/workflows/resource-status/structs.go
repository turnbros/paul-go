package resource_status

type StatusRequest struct {
	ResourceScope string `json:"resource_scope"`
	ResourceType  string `json:"resource_type"`
}
