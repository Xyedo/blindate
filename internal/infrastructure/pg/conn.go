package pg

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/xyedo/blindate/internal/infrastructure"
)

type Querier interface {
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults

	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

var (
	pool *pgxpool.Pool
	once sync.Once
)

func Transaction(ctx context.Context, option pgx.TxOptions, cb func(tx Querier) error) error {
	tx, err := GetPool(ctx).BeginTx(ctx, option)

	if err != nil {
		return err
	}

	err = cb(tx)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return fmt.Errorf("cannot rollback with error %v:%v", txErr, err)
		}

		return err

	}

	return tx.Commit(ctx)

}

func GetPool(ctx context.Context) *pgxpool.Pool {
	once.Do(func() {
		_, _ = InitPool(ctx)
	})

	return pool
}

func InitPool(ctx context.Context) (*pgxpool.Pool, error) {
	var err error
	once.Do(func() {
		config := infrastructure.Config.DbConf

		connStr := fmt.Sprintf(
			"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
		)

		if conn, pgErr := pgxpool.New(ctx, connStr); pgErr != nil {
			err = pgErr
		} else {
			pool = conn
		}
	})
	if err != nil {
		return nil, err
	}

	return pool, nil

}
