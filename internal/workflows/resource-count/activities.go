package resource_count

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"paul/internal/util"
)

func CountAll(ctx context.Context, countRequest CountRequest) (string, error) {

	podList := util.ListKubePods("", metav1.ListOptions{})
	count := len(podList.Items)
	fmt.Println("Pod Count is: ", len(podList.Items))
	return fmt.Sprintf("The answer is %v", count), nil
}

func CountPods(ctx context.Context, countRequest CountRequest) (string, error) {
	podList := util.ListKubePods(countRequest.ResourceScope, metav1.ListOptions{})
	count := len(podList.Items)

	if countRequest.ResourceScope != "" {
		return fmt.Sprintf("The answer is %v pods in %v", count, countRequest.ResourceScope), nil
	}
	return fmt.Sprintf("The answer is %v pods in total", count), nil
}

func CountServices(ctx context.Context, countRequest CountRequest) (string, error) {
	serviceList := util.ListKubeServices(countRequest.ResourceScope, metav1.ListOptions{})
	count := len(serviceList.Items)

	if countRequest.ResourceScope != "" {
		return fmt.Sprintf("The answer is %v pods in %v", count, countRequest.ResourceScope), nil
	}
	return fmt.Sprintf("The answer is %v pods in total", count), nil
}
