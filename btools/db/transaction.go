package db

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/seemyown/backend-toolkit/btools/exc"
	"time"
)

var trxLogger = log.NewSubLogger("trx")

type Transaction interface {
	Exec(ctx context.Context, fn func(tx *sqlx.Tx) error) error
}

type trx struct {
	db *sqlx.DB
}

func NewTrx(db *Database) Transaction {
	return &trx{db: db.DB}
}

func (t *trx) Exec(ctx context.Context, fn func(tx *sqlx.Tx) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		log.Error(err, "Error starting transaction")
		return exc.RepositoryError("transaction_begin_error")
	}
	startTime := time.Now()
	if err := fn(tx); err != nil {
		log.Error(err, "Error executing transaction. Rollback...")
		_ = tx.Rollback()
		return err
	}
	log.Info(fmt.Sprintf("Transaction finished in %f seconds", time.Since(startTime).Seconds()))
	if err := tx.Commit(); err != nil {
		log.Error(err, "Error committing transaction")
		return exc.RepositoryError("transaction_commit_error")
	}
	return nil
}
