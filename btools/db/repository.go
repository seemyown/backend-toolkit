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
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context, offset, limit int64, args ...interface{}) ([]*T, error)
	Search(ctx context.Context, offset, limit int64, args ...interface{}) ([]*T, error)
}

type BaseRepository[T any] struct {
	Db  *sqlx.DB
	Trx *Trx
}

func (r *BaseRepository[T]) Get(ctx context.Context, id int64) (*T, error) {
	return nil, exc.RepositoryError("Not implemented")
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return exc.RepositoryError("Not implemented")
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id int64) error {
	return exc.RepositoryError("Not implemented")
}

func (r *BaseRepository[T]) GetAll(ctx context.Context, offset, limit int64, args ...interface{}) ([]*T, error) {
	return make([]*T, 0), exc.RepositoryError("Not implemented")
}

func NewBaseRepository[T any](conn *Database) *BaseRepository[T] {
	return &BaseRepository[T]{
		Db:  conn.DB,
		Trx: NewTrx(conn.DB),
	}
}

func (r *BaseRepository[T]) GetSingleRow(ctx context.Context, query string, args ...interface{}) (*T, error) {
	stmt, err := r.Db.PreparexContext(ctx, query)
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

func (r *BaseRepository[T]) SelectManyRow(ctx context.Context, query string, args ...interface{}) ([]*T, error) {
	stmt, err := r.Db.PreparexContext(ctx, query)
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

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return exc.RepositoryError("Not implemented")
}
