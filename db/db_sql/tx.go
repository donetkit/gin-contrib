package db_sql

import (
	"context"
	"database/sql/driver"

	"go.opentelemetry.io/otel/trace"
)

type sqlTx struct {
	ctx     context.Context
	tx      driver.Tx
	instrum *dbInstrum
}

var _ driver.Tx = (*sqlTx)(nil)

func newTx(ctx context.Context, tx driver.Tx, instrum *dbInstrum) *sqlTx {
	return &sqlTx{
		ctx:     ctx,
		tx:      tx,
		instrum: instrum,
	}
}

func (tx *sqlTx) Commit() error {
	return tx.instrum.withSpan(tx.ctx, "tx.Commit", "",
		func(ctx context.Context, span trace.Span) error {
			return tx.tx.Commit()
		})
}

func (tx *sqlTx) Rollback() error {
	return tx.instrum.withSpan(tx.ctx, "tx.Rollback", "",
		func(ctx context.Context, span trace.Span) error {
			return tx.tx.Rollback()
		})
}
