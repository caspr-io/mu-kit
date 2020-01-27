package config

import "k8s.io/client-go/rest"

import "github.com/caspr-io/mu-kit/util"

type Loader interface {
	Load() (*rest.Config, error)
}

func fallback(err error, loader Loader) (*rest.Config, error) {
	collector := &util.ErrorCollector{}
	collector.Collect(err)

	if loader != nil {
		c, e := loader.Load()
		if e != nil {
			collector.Collect(e)
			return nil, collector
		}

		return c, nil
	}

	return nil, collector
}
