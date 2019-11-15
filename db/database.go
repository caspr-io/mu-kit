package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v9"
)

type PostgreSQLConfig struct {
	Host     string `split_words:"true" required:"true"`
	Port     int    `split_words:"true" required:"true" default:"5432"`
	User     string `split_words:"true" required:"true"`
	Password string `split_words:"true" required:"true"`
	Database string `split_words:"true" required:"true"`
	PoolSize int    `split_words:"true" required:"true" default:"10"`
}

func ConnectToPostgreSQL(config *PostgreSQLConfig) *pg.DB {
	return pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		User:     config.User,
		Password: config.Password,
		PoolSize: config.PoolSize,
	})
}

func AsDatabaseSQL(pgDB *pg.DB) (*sql.DB, error) {
	opts := pgDB.Options()

	sslMode := "disable"
	if opts.TLSConfig != nil {
		sslMode = "require"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		strings.Split(opts.Addr, ":")[0], strings.Split(opts.Addr, ":")[1], opts.User, opts.Password, opts.Database, sslMode)

	return sql.Open("postgres", psqlInfo)
}
