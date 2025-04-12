package repository

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/flags"
	internalErrors "github.com/ramil063/firstgodiplom/internal/errors"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

type Repository struct {
	Pool *pgxpool.Pool
}

// DataBaser основные команды для работы с бд
type DataBaser interface {
	ExecContext(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	QueryRowContext(ctx context.Context, query string, args ...any) pgx.Row
	QueryContext(ctx context.Context, query string, args ...any) (pgx.Rows, error)
	Open() (*pgxpool.Pool, error)
	PingContext(ctx context.Context) error
	SetPool() error
}

var DBRepository Repository

func (dbr *Repository) ExecContext(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	result, err := dbr.Pool.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (dbr *Repository) QueryRowContext(ctx context.Context, query string, args ...any) pgx.Row {
	row := dbr.Pool.QueryRow(ctx, query, args...)
	return row
}

func (dbr *Repository) QueryContext(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	rows, err := dbr.Pool.Query(ctx, query, args...)

	if err != nil {
		return nil, internalErrors.NewDBError(err)
	}
	return rows, nil
}

// Open открыть соединение с бд
func (dbr *Repository) Open() (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(flags.DatabaseURI)
	if err != nil {
		return nil, internalErrors.NewDBError(err)
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// do something with every new connection
		return nil
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)

	if err != nil {
		return nil, internalErrors.NewDBError(err)
	}
	return pool, nil
}

// PingContext проверить соединение с бд
func (dbr *Repository) PingContext(ctx context.Context) error {
	err := dbr.Pool.Ping(ctx)
	if err != nil {
		logger.WriteErrorLog(err.Error())
	}
	return err
}

// SetPool установить поле Pool для работы с бд
func (dbr *Repository) SetPool() error {
	pool, err := dbr.Open()
	if err != nil {
		logger.WriteErrorLog(err.Error())
		return err
	}
	dbr.Pool = pool
	return nil
}

func NewRepository() (*Repository, error) {
	rep := &Repository{}
	err := rep.SetPool()
	return rep, err
}
