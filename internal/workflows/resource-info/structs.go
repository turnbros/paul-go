package resource_info

type InfoRequest struct {
	ResourceName  string `json:"resource_name"`
	ResourceScope string `json:"resource_scope"`
	ResourceType  string `json:"resource_type"`
}
