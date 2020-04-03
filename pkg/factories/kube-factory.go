package factories

import (
	"flag"
	"path/filepath"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	lazyConfig sync.Once
	lazyClient sync.Once
	kubeClient *kubernetes.Clientset
	kubeConfig *rest.Config
)

// CreateKubeConfig creates and loads kubeconfig from file or in cluster
func CreateKubeConfig() *rest.Config {
	lazyConfig.Do(func() {
		var kconfig *string

		if home := homedir.HomeDir(); home != "" {
			path := filepath.Join(home, ".kube", "config")
			kconfig = flag.String("kubeconfig", path, "absolute path to teh kubeconfig file")
		} else {
			kconfig = flag.String("kubeconfig", "", "absolute path to kubeconfig file")
		}

		flag.Parse()

		conf, err := rest.InClusterConfig()
		if err != nil {
			// try cmdline
			conf, err = clientcmd.BuildConfigFromFlags("", *kconfig)
			if err != nil {
				panic(err)
			}
		}

		kubeConfig = conf
	})

	return kubeConfig
}

// CreateKubeClient creates a in cluster Kubernetes client
func CreateKubeClient() *kubernetes.Clientset {
	lazyClient.Do(func() {
		config := CreateKubeConfig()

		client, err := kubernetes.NewForConfig(config)
		if err != nil {
			panic(err)
		}

		kubeClient = client
	})
	return kubeClient
}
