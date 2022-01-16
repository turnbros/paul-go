package activities

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"paul-go/internal/util"
	rs "paul-go/internal/workflows/resource-info/util"
)

func CountNamespaces(ctx context.Context, resourceRequest rs.ResourceRequest) (string, error) {
	namespaceList := util.ListKubeNamespaces(metav1.ListOptions{})
	count := len(namespaceList.Items)

	if resourceRequest.ResourceKubeNamespace != "" {
		return fmt.Sprintf("There are %v namespaces in %v", count, resourceRequest.ResourceKubeNamespace), nil
	}
	return fmt.Sprintf("There are %v namespaces total", count), nil
}

func CountPods(ctx context.Context, resourceRequest rs.ResourceRequest) (string, error) {
	podList := util.ListKubePods(resourceRequest.ResourceKubeNamespace, metav1.ListOptions{})
	count := len(podList.Items)

	if resourceRequest.ResourceKubeNamespace != "" {
		return fmt.Sprintf("There are %v pods in %v", count, resourceRequest.ResourceKubeNamespace), nil
	}
	return fmt.Sprintf("There are %v pods total", count), nil
}

func CountServices(ctx context.Context, resourceRequest rs.ResourceRequest) (string, error) {
	serviceList := util.ListKubeServices(resourceRequest.ResourceKubeNamespace, metav1.ListOptions{})
	count := len(serviceList.Items)

	if resourceRequest.ResourceKubeNamespace != "" {
		return fmt.Sprintf("There are %v services in %v", count, resourceRequest.ResourceKubeNamespace), nil
	}
	return fmt.Sprintf("There are %v services total", count), nil
}
