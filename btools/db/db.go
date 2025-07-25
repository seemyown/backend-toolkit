package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/seemyown/backend-toolkit/btools/exc"
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

func SelectOne[T any](db *sqlx.DB, ctx context.Context, query string, args ...interface{}) (*T, error) {
	var result T
	if err := db.GetContext(ctx, &result, query, args...); err != nil {
		Logger.Error(err, "failed to execute query %s, %v", query, args)
		return nil, WrapError(err)
	}
	return &result, nil
}

func SelectMany[T any](db *sqlx.DB, ctx context.Context, query string, args ...interface{}) ([]*T, error) {
	var result []*T

	if err := db.SelectContext(ctx, &result, query, args...); err != nil {
		Logger.Error(err, "failed to execute query %s, %v", query, args)
		return nil, WrapError(err)
	}
	return result, nil
}

func SelectIn[T any](ctx context.Context, db *sqlx.DB, dest *[]T, query string, inArgs any) error {
	q, args, err := sqlx.In(query, inArgs)
	if err != nil {
		return exc.RepositoryError(err.Error())
	}

	q = db.Rebind(q)

	if err := db.SelectContext(ctx, dest, q, args...); err != nil {
		return WrapError(err)
	}

	return nil
}

func PrepareAndExec[T any](ctx context.Context, tx *sqlx.Tx, query string, items []*T, execFn func(*sqlx.Stmt, *T) error) error {
	stmt, err := tx.PreparexContext(ctx, query)
	if err != nil {

		return exc.RepositoryError(err.Error())
	}
	defer func() { _ = stmt.Close() }()

	for _, item := range items {
		if err := execFn(stmt, item); err != nil {
			return WrapError(err)
		}
	}
	return nil
}
