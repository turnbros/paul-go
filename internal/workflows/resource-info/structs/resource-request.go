package resource_info_structs

type ResourceRequest struct {
	RequestType            string `json:"request_type"`
	ResourceName           string `json:"resource_name"`
	ResourceType           string `json:"resource_type"`
	ResourceKubeNamespace  string `json:"resource_kube_namespace"`
	ResourceKubeService    string `json:"resource_kube_service"`
	ResourceKubeDeployment string `json:"resource_kube_deployment"`
	ResourceKubePod        string `json:"resource_kube_pod"`
	ResourceKubeNode       string `json:"resource_kube_node"`
}
