package resources_count

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"paul/internal/util"
	"strconv"
)

func CountAll(ctx context.Context) (string, error) {

	podList := util.ListKubePods("", metav1.ListOptions{})
	count := len(podList.Items)
	fmt.Println("Pod Count is: ", len(podList.Items))
	return strconv.Itoa(count), nil
}
