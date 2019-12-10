package config

import (
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type FileLoader struct {
	File     string
	Fallback Loader
}

func (l *FileLoader) Load() (*rest.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", l.File)
	if err != nil {
		log.Error().Err(err).Str("kube-config", l.File).Msg("Cannot load Kubernetes configuration")
		return fallback(err, l.Fallback)
	}

	return config, nil
}
