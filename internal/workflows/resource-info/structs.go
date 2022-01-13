package resource_info

type InfoRequest struct {
	RequestType   string `json:"request_type"`
	ResourceName  string `json:"resource_name"`
	ResourceScope string `json:"resource_scope"`
	ResourceType  string `json:"resource_type"`
	//ResourceType  string `json:"resource_type"`
	// resource_kube_namespace
	// resource_kube_service
	// resource_kube_deployment
	// resource_kube_pod
	// resource_kube_node
}
