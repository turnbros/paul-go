package resource_info

import (
	"bytes"
	"context"
	"fmt"
	"github.com/olekukonko/tablewriter"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"paul/internal/util"
	"time"
)

func ListPods(ctx context.Context, infoRequest InfoRequest) (string, error) {
	podList := util.ListKubePods(infoRequest.ResourceScope, metav1.ListOptions{}).Items

	var data [][]string
	for i := 1; i < len(podList); i++ {
		data = append(data, []string{
			podList[i].Name,
			podList[i].Namespace,
			podList[i].Status.String(),
			time.Since(podList[i].CreationTimestamp.Time).String(),
		})
	}
	renderedPodList := renderTable(data)

	return fmt.Sprintf("Here you go!\n```%v```", renderedPodList), nil
}

func renderTable(tableData [][]string) string {
	byteBuffer := new(bytes.Buffer)
	table := tablewriter.NewWriter(byteBuffer)
	table.SetHeader([]string{"NAME", "NAMESPACE", "READY", "STATUS", "AGE"})

	for _, v := range tableData {
		table.Append(v)
	}
	table.Render()
	return byteBuffer.String()
}
