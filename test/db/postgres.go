package db

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/caspr-io/mu-kit/db"
	database "github.com/caspr-io/mu-kit/db"
	"github.com/caspr-io/mu-kit/test/docker"
	"github.com/rs/zerolog/log"

	"github.com/go-pg/pg/v9"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // this is a comment jutifying the import ;)
)

func Postgres(dckr *docker.Docker) (*pg.DB, error) {
	log.Logger.Info().Msg("Starting Postgres Docker image")

	c, err := dckr.RunContainer("postgres", "12", []string{"POSTGRES_PASSWORD=secret"})
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=localhost port=%s user=postgres password=secret dbname=postgres sslmode=disable", c.GetPort("5432/tcp"))

	if err := c.WaitForRunning(func() error {
		var err error
		db, err := sql.Open("postgres", psqlInfo)
		defer db.Close() //nolint:staticcheck

		if err != nil {
			return err
		}

		return db.Ping()
	}); err != nil {
		return nil, err
	}

	port, err := strconv.Atoi(c.GetPort("5432/tcp"))
	if err != nil {
		return nil, err
	}

	return db.ConnectToPostgreSQL(&database.PostgreSQLConfig{
		Host:     "localhost",
		Port:     port,
		User:     "postgres",
		Password: "secret",
		Database: "postgres",
	}), nil
}

func Migrate(pgDB *pg.DB, migrations string) error {
	db, err := database.AsDatabaseSQL(pgDB)
	if err != nil {
		return err
	}

	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		migrations,
		"postgres", driver)
	if err != nil {
		return err
	}

	defer m.Close()

	return m.Up()
}
