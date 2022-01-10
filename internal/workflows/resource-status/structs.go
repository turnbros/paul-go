package resource_status

type StatusRequest struct {
	// TODO: Look into way to add a scope to this
	// ResourceScope string `json:"resource_scope"`
	ResourceType string `json:"resource_type"`
}
