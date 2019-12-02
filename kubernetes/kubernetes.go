package kubernetes

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8s{config, client}, nil
}

func getKubernetesConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Info().
			Err(err).
			Msg("Cannot load in-cluster k8s configuration, now trying to load config from ~/.kube/config")

		home, err := homedir.Dir()
		if err != nil {
			log.Error().
				Err(err).
				Msg("Cannot determine home directory")

			return nil, err
		}

		kubeconfig := filepath.Join(home, ".kube", "config")

		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Error().
				Err(err).
				Str("kube-config", kubeconfig).
				Msg("Cannot load Kubernetes configuration")

			return nil, err
		}
	}

	return config, nil
}
