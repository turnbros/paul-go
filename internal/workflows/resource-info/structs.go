package resource_info

type InfoRequest struct {
	ResourceScope string `json:"resource_scope"`
	ResourceType  string `json:"resource_type"`
}
