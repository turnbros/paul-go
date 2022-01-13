package activities

import (
	"context"
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"paul/internal/util"
	util2 "paul/internal/workflows/resource-info/util"
	"time"
)

func ListPods(ctx context.Context, resourceRequest util2.ResourceRequest) (string, error) {
	podList := util.ListKubePods(resourceRequest.ResourceKubeNamespace, metav1.ListOptions{}).Items

	var data [][]string
	for i := 1; i < len(podList); i++ {
		data = append(data, []string{
			podList[i].Name,
			podList[i].Namespace,
			podList[i].Status.Reason,
			time.Since(podList[i].CreationTimestamp.Time).String(),
		})
	}
	renderedPodList := util2.RenderTable(data)

	return fmt.Sprintf("Here you go!\n%v", renderedPodList), nil
}

func GetPodInfo(ctx context.Context, infoRequest util2.ResourceRequest) (string, error) {

	return fmt.Sprintf("Here you go!\n```%v```", ""), nil
}
