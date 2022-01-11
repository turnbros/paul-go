package util

import (
	"context"
	"errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func getKubeConfig() *rest.Config {
	var config *rest.Config

	kubeConfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	if _, err := os.Stat(kubeConfigPath); err == nil {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			panic(err)
		}

	} else if errors.Is(err, os.ErrNotExist) {
		/*svcToken, readErr := os.ReadFile("/run/secrets/kubernetes.io/serviceaccount/token")
		if readErr != nil {
			panic(readErr)
		}*/
		config = &rest.Config{
			Host:            "https://" + os.Getenv("KUBERNETES_SERVICE_HOST") + ":" + os.Getenv("KUBERNETES_SERVICE_PORT"),
			BearerTokenFile: "/run/secrets/kubernetes.io/serviceaccount/token",
			TLSClientConfig: rest.TLSClientConfig{
				Insecure: false,
				CAFile:   "/run/secrets/kubernetes.io/serviceaccount/ca.crt",
			},
		}
	}

	return config
}

func GetKubeClient() *kubernetes.Clientset {
	config := getKubeConfig()
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func ListKubeNamespaces(options metav1.ListOptions) *v1.NamespaceList {
	kubeClient := GetKubeClient()
	namespaceList, err := kubeClient.CoreV1().Namespaces().List(context.TODO(), options)
	if err != nil {
		panic(err.Error())
	}
	return namespaceList
}

func ListKubePods(namespace string, options metav1.ListOptions) *v1.PodList {
	kubeClient := GetKubeClient()
	podList, err := kubeClient.CoreV1().Pods(namespace).List(context.TODO(), options)
	if err != nil {
		panic(err.Error())
	}
	return podList
}

func ListKubeServices(namespace string, options metav1.ListOptions) *v1.ServiceList {
	kubeClient := GetKubeClient()
	serviceList, err := kubeClient.CoreV1().Services(namespace).List(context.TODO(), options)
	if err != nil {
		panic(err.Error())
	}
	return serviceList
}

func GetConfigMap(namespace string, configName string) *v1.ConfigMap {
	ctx := context.Background()
	kubeClient := GetKubeClient()
	configMap, err := kubeClient.CoreV1().ConfigMaps(namespace).Get(ctx, configName, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	return configMap
}

func GetConfigMapData(namespace string, configName string) map[string]string {
	return GetConfigMap(namespace, configName).Data
}
