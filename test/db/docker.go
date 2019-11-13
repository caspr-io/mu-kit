package db

import (
	"database/sql"
	"fmt"

	"github.com/ory/dockertest/v3"
	"github.com/rs/zerolog/log"
)

func StartPostgresContainer() (*sql.DB, error) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Logger.Error().Err(err).Msg("Could not connect to docker")
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.Run("postgres", "12", []string{"POSTGRES_PASSWORD=secret"})
	if err != nil {
		log.Logger.Error().Err(err).Msg("Could not start resource")
	}

	psqlInfo := fmt.Sprintf("host=localhost port=%s user=postgres password=secret dbname=postgres sslmode=disable", resource.GetPort("5432/tcp"))

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		var err error
		db, err := sql.Open("postgres", psqlInfo)
		defer db.Close() //nolint:staticcheck

		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		log.Logger.Error().Err(err).Msg("Could not connect to dockerized postgres")
	}

	return sql.Open("postgres", psqlInfo)
}
