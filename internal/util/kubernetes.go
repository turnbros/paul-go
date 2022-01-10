package util

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func GetKubeClient() *kubernetes.Clientset {

	// var kubeconfig *string

	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	/*if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()*/

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

func ListKubePods(namespace string, options metav1.ListOptions) *v1.PodList {
	kubeClient := GetKubeClient()
	podList, err := kubeClient.CoreV1().Pods("").List(context.TODO(), options)
	if err != nil {
		panic(err.Error())
	}
	return podList
}
