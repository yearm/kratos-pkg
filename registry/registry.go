package registry

import (
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/sirupsen/logrus"
	"github.com/yearm/kratos-pkg/config/env"
	kuberegistry "github.com/yearm/kratos-pkg/registry/kubernetes"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// NewRegistrar ...
func NewRegistrar() registry.Registrar {
	if !env.IsLocal() {
		reg, err := NewKubeRegistry()
		if err != nil {
			logrus.Panicln("new registrar error:", err)
		}
		return reg
	}
	return nil
}

// NewDiscovery ...
func NewDiscovery() registry.Discovery {
	if !env.IsLocal() {
		discovery, err := NewKubeDiscovery()
		if err != nil {
			logrus.Panicln("new discovery error:", err)
		}
		return discovery
	}
	return nil
}

// NewKubeRegistry ...
func NewKubeRegistry() (registry.Registrar, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	reg := kuberegistry.NewRegistry(clientSet)
	reg.Start()
	return reg, nil
}

// NewKubeDiscovery ...
func NewKubeDiscovery() (registry.Discovery, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	discovery := kuberegistry.NewRegistry(clientSet)
	discovery.Start()
	return discovery, nil
}
