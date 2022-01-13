package activities

import (
	"context"
	"fmt"
	util2 "paul/internal/workflows/resource-info/util"
)

func GetPodInfo(ctx context.Context, infoRequest util2.ResourceRequest) (string, error) {

	return fmt.Sprintf("Here you go!\n```%v```", ""), nil
}
