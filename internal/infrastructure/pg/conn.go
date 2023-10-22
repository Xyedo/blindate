package pg

import (
	"context"
	"fmt"

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

func Transaction(ctx context.Context, option pgx.TxOptions, cb func(tx Querier) error) error {
	pool, err := GetPool(ctx)
	if err != nil {
		return err
	}

	defer pool.Close()

	tx, err := pool.BeginTx(ctx, option)

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

func GetPool(ctx context.Context) (*pgxpool.Pool, error) {
	config := infrastructure.Config.DbConf

	connStr := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	conn, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
