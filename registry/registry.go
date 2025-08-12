package registry

import (
	"github.com/go-kratos/kratos/v2/registry"
	kuberegistry "github.com/yearm/kratos-pkg/registry/kubernetes"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewKubeRegistry creates and initializes a Kubernetes-based service registry.
func NewKubeRegistry() (registry.Registrar, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	reg := kuberegistry.NewRegistry(clientSet, "")
	reg.Start()
	return reg, nil
}

// NewKubeDiscovery creates a Kubernetes service discovery client for service resolution.
func NewKubeDiscovery() (registry.Discovery, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	discovery := kuberegistry.NewRegistry(clientSet, "")
	discovery.Start()
	return discovery, nil
}
