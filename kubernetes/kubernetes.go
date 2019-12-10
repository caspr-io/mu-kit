package kubernetes

import (
	"github.com/caspr-io/mu-kit/kubernetes/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8s struct {
	Config *rest.Config
	*kubernetes.Clientset
}

func ConnectToKubernetes() (*K8s, error) {
	config, err := getKubernetesConfig()
	if err != nil {
		return nil, err
	}

	return ConnectUsingConfig(config)
}

func ConnectUsingConfig(config *rest.Config) (*K8s, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8s{config, client}, nil
}

func getKubernetesConfig() (*rest.Config, error) {
	return (&config.InClusterConfigLoader{
		Fallback: &config.HomeDirLoader{},
	}).Load()
}
