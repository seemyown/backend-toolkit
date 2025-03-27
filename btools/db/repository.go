package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/seemyown/backend-toolkit/btools/exc"
	"github.com/seemyown/backend-toolkit/btools/logging"
)

var Logger = logging.New(logging.Config{
	FileName: "repository",
	Name:     "base",
})

type Repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Get(ctx context.Context, id int64) (*T, error)
}

type BaseRepository[T any] struct {
	db  *sqlx.DB
	trx *Trx
}

func NewBaseRepository[T any](conn *sqlx.DB) *BaseRepository[T] {
	return &BaseRepository[T]{
		db:  conn,
		trx: NewTrx(conn),
	}
}

func prepareContextSindgleRow[T any](ctx context.Context, db *sqlx.DB, query string, args ...interface{}) (*T, error) {
	stmt, err := db.PreparexContext(ctx, query)
	if err != nil {
		Logger.Error(err, "failed to prepare statement")
		return nil, exc.RepositoryError(err.Error())
	}
	defer func() { _ = stmt.Close() }()

	var result T
	if err := stmt.GetContext(ctx, &result, args...); err != nil {
		Logger.Error(err, "failed to execute query %s, %p", query, args)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, exc.RepositoryError(err.Error())
	}
	return &result, nil
}

func prepareContextManyRow[T any](ctx context.Context, db *sqlx.DB, query string, args ...interface{}) ([]*T, error) {
	stmt, err := db.PreparexContext(ctx, query)
	if err != nil {
		Logger.Error(err, "failed to prepare statement")
		return nil, exc.RepositoryError(err.Error())
	}
	defer func() { _ = stmt.Close() }()

	var result []*T
	if err := stmt.SelectContext(ctx, &result, args...); err != nil {
		Logger.Error(err, "failed to execute query %s, %p", query, args)
		if errors.Is(err, sql.ErrNoRows) {
			return make([]*T, 0), nil
		}
		return nil, exc.RepositoryError(err.Error())
	}
	return result, nil
}
