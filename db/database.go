package database

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v9"
)

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
