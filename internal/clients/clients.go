package clients

import (
	"fmt"

	configclient "github.com/openshift/client-go/config/clientset/versioned"
	imageclient "github.com/openshift/client-go/image/clientset/versioned"
	projectclient "github.com/openshift/client-go/project/clientset/versioned"
	routeclient "github.com/openshift/client-go/route/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Clients struct {
	Kubernetes *kubernetes.Clientset
	Config     *configclient.Clientset
	Route      *routeclient.Clientset
	Image      *imageclient.Clientset
	Project    *projectclient.Clientset
	RestConfig *rest.Config
}

func NewClients() (*Clients, error) {
	config, err := getKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	configClient, err := configclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create config client: %w", err)
	}

	routeClient, err := routeclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create route client: %w", err)
	}

	imageClient, err := imageclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create image client: %w", err)
	}

	projectClient, err := projectclient.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create project client: %w", err)
	}

	return &Clients{
		Kubernetes: kubeClient,
		Config:     configClient,
		Route:      routeClient,
		Image:      imageClient,
		Project:    projectClient,
		RestConfig: config,
	}, nil
}
