package pg

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"github.com/pkg/errors"

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
	conn, err := GetConnectionPool(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	tx, err := conn.BeginTx(ctx, option)
	if err != nil {
		return err
	}

	err = cb(tx)
	if err != nil {
		if txErr := tx.Rollback(ctx); txErr != nil {
			return fmt.Errorf("err %w but cannot rollback with error:%w", err, txErr)
		}

		return err

	}

	return tx.Commit(ctx)
}

var ErrMaxRetry = errors.New("retry limit exceeded")

const maxRetry = 5

func TransactionWithRetry(ctx context.Context, txOpt pgx.TxOptions, cb func(tx Querier) error) error {
	for i := 0; i < maxRetry; i++ {
		conn, err := GetConnectionPool(ctx)
		if err != nil {
			return err
		}

		tx, err := conn.BeginTx(ctx, txOpt)
		if err != nil {
			conn.Release()
			return err
		}

		err = cb(tx)
		if err != nil {

			if rbErr := tx.Rollback(ctx); rbErr != nil {
				conn.Release()
				return fmt.Errorf("cannot rollback %w: %w", rbErr, err)
			}
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "40001" {
				conn.Release()
				continue

			}
			conn.Release()
			return err

		}

		err = tx.Commit(ctx)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "40001" {
				conn.Release()
				continue
			}

			conn.Release()
			return err
		}

		conn.Release()
		return nil
	}

	return ErrMaxRetry

}

func GetConnectionPool(ctx context.Context) (*pgxpool.Conn, error) {
	if pool == nil {
		once.Do(func() {
			_, _ = InitPool(ctx)
		})
	}

	return pool.Acquire(ctx)
}

func InitPool(ctx context.Context) (*pgxpool.Pool, error) {
	config := infrastructure.Config.Postgre

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable pool_max_conns=%d",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		runtime.NumCPU()*4,
	)

	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(ctx); err != nil {
		return nil, err
	}

	pool = db

	return db, nil

}
