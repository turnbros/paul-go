package resource_info

type InfoRequest struct {
	// TODO: Look into way to add a scope to this
	// ResourceScope string `json:"resource_scope"`
	ResourceType string `json:"resource_type"`
}
