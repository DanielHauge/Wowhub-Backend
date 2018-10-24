package Postgres

import (
	log "../Utility/Logrus"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
)

var db *sql.DB

func init() {
	connectionPool, err := sql.Open("postgres", "postgres://"+os.Getenv("DBPASS")+"@localhost/postgres?sslmode=disable")
	if err != nil {
		log.WithLocation().WithError(err).Fatal("Could not establish connection pool for postgres")
	}
	e := connectionPool.Ping()
	if e != nil {
		log.WithLocation().WithError(e).Fatal("Could not establish connection to postgres")
	}

	connectionPool.SetMaxIdleConns(1)
	connectionPool.SetMaxOpenConns(10)
	db = connectionPool

}