package util

import (
	"bytes"
	"github.com/olekukonko/tablewriter"
)

func RenderTable(tableData [][]string) string {
	byteBuffer := new(bytes.Buffer)
	table := tablewriter.NewWriter(byteBuffer)
	table.SetHeader([]string{"NAME", "NAMESPACE", "STATUS", "AGE"})

	for _, v := range tableData {
		table.Append(v)
	}
	table.Render()
	return byteBuffer.String()
}
