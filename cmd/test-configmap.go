package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
	"paul/internal/util"
)

func main() {
	log.Println("here we go")
	log.Println("the kind: ", util.ListKubePods("", metav1.ListOptions{}).Kind)

	temporalConfig := util.GetTemporalConfig()

	log.Println(temporalConfig)
}
