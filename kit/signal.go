package kit

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
)

type Closer interface {
	Close() error
}

func SignalsHandler(closer Closer, logger zerolog.Logger) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		logger.Info().Interface("signal", sig).Msg("Received signal, closing")

		err := closer.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Close failed")
		}
	}()

	return nil
}
