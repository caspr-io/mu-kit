package config

import (
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
)

type InClusterConfigLoader struct {
	Fallback Loader
}

func (l *InClusterConfigLoader) Load() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Info().
			Err(err).
			Msg("Cannot load in-cluster k8s configuration, falling back")
		return fallback(err, l.Fallback)
	}
	return config, nil
}
