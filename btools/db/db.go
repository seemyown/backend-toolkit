package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/seemyown/backend-toolkit/btools/logging"
)

var log = logging.New(logging.Config{
	FileName: "repository",
	Name:     "db",
})

type Database struct {
	DB *sqlx.DB
}

func NewDatabase(dsn string) *Database {
	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Error(err, "error connecting to database")
		panic(err)
	}
	return &Database{conn}
}
