package k8s

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Config struct {
	Kubeconfig string
	InCluster  bool
}

func BuildRestConfig(cfg Config) (*rest.Config, error) {
	if cfg.InCluster {
		return rest.InClusterConfig()
	}

	if cfg.Kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", cfg.Kubeconfig)
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	overrides := &clientcmd.ConfigOverrides{}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides).ClientConfig()
}

func NewClients(restCfg *rest.Config) (*kubernetes.Clientset, dynamic.Interface, error) {
	kubeClient, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return nil, nil, err
	}

	dynClient, err := dynamic.NewForConfig(restCfg)
	if err != nil {
		return nil, nil, err
	}

	return kubeClient, dynClient, nil
}
