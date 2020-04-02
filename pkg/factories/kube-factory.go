package factories

import (
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	lazy       sync.Once
	kubeClient *kubernetes.Clientset
)

// CreateKubeClient creates a in cluster Kubernetes client
func CreateKubeClient() *kubernetes.Clientset {

	lazy.Do(func() {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err)
		}

		kubeClient, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}
	})
	return kubeClient
}
