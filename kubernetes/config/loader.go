package config

import "k8s.io/client-go/rest"

type Loader interface {
	Load() (*rest.Config, error)
}

func fallback(err error, loader Loader) (*rest.Config, error) {
	if loader != nil {
		return loader.Load()
	} else {
		return nil, err
	}
}
