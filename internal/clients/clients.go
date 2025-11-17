package clients

import (
	"fmt"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/rest"
)

type Clients struct {
	Kubernetes *kubernetes.Clientset
	Dynamic    dynamic.Interface
	RestConfig *rest.Config
}

func NewClients() (*Clients, error) {
	cfg, err := getKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	kube, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	dyn, err := dynamic.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &Clients{
		Kubernetes: kube,
		Dynamic:    dyn,
		RestConfig: cfg,
	}, nil
}
