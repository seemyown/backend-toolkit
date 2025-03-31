package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/seemyown/backend-toolkit/btools/logging"
)

var log = logging.New(logging.Config{
	FileName: "repository",
	Name:     "db",
})

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Params   map[string]string
}

func (c *Config) String() string {
	baseConnString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database,
	)

	for key, value := range c.Params {
		baseConnString += fmt.Sprintf(" %s=%s", key, value)
	}

	return baseConnString
}

type Database struct {
	DB *sqlx.DB
}

func NewDatabase(cfg *Config) *Database {
	conn, err := sqlx.Connect("postgres", cfg.String())
	if err != nil {
		log.Error(err, "error connecting to database")
		panic(err)
	}
	return &Database{conn}
}
