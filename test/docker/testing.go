package docker

import (
	"github.com/caspr-io/mu-kit/util"
	"github.com/rs/zerolog/log"
	"io"
	"testing"
)

func RunTestsWithDockerContainers(m *testing.M, fs ...func(*Docker) (io.Closer, error)) (int, error) {
	closer := new(util.MultiCloser)
	dckr, err := StartDocker()

	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Could not start Docker")
	}

	closer.Add(dckr)

	errorCollector := new(util.ErrorCollector)

	for _, f := range fs {
		c, err := f(dckr)
		closer.Add(c) // Even if an error is returned, we could have an open connection which needs closing

		if err != nil {
			errorCollector.Collect(err)
			break
		}
	}

	if errorCollector.HasErrors() {
		if err := closer.Close(); err != nil {
			errorCollector.Collect(err)
		}

		return -1, errorCollector
	}

	code := m.Run()
	log.Logger.Info().Int("code", code).Msg("Ran tests")

	if err := closer.Close(); err != nil {
		return -1, err
	}

	return code, nil
}
