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

func CountPods(ctx context.Context) (string, error) {

	podList := util.ListKubePods("", metav1.ListOptions{})
	count := len(podList.Items)
	fmt.Println("Pod Count is: ", len(podList.Items))
	return fmt.Sprintf("The answer is %v", count), nil
}
