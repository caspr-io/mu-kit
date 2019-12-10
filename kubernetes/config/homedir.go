package config

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"k8s.io/client-go/rest"
)

type HomeDirLoader struct {
	Fallback Loader
}

func (l *HomeDirLoader) Load() (*rest.Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		log.Error().Err(err).Msg("Cannot determine home directory")
		return fallback(err, l.Fallback)
	}

	kubeconfig := filepath.Join(home, ".kube", "config")

	return (&FileLoader{kubeconfig, l.Fallback}).Load()
}
