package store

import (
	"context"

	"adtech.simple/internal/pkg/dbquery"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	dbPool *pgxpool.Pool
}

func NewStorage(dbPool *pgxpool.Pool) *Storage {
	return &Storage{
		dbPool: dbPool,
	}
}

func (s *Storage) WithTx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx, queries *dbquery.Queries) error) (err error) {
	tx, err := s.dbPool.Begin(ctx)
	if err != nil {
		return
	}

	dbQuerier := dbquery.New(s.dbPool).WithTx(tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
		if err != nil {
			_ = tx.Rollback(ctx)
			return
		}

		err = tx.Commit(ctx)
	}()

	err = fn(ctx, tx, dbQuerier)

	return
}
