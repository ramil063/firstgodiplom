package dml

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/flags"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
	internalErrors "github.com/ramil063/firstgodiplom/internal/errors"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

type Repository struct {
	Pool *pgxpool.Pool
}

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

func (dbr *Repository) PingContext(ctx context.Context) error {
	err := dbr.Pool.Ping(ctx)
	if err != nil {
		logger.WriteErrorLog(err.Error())
	}
	return err
}

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

func AddUser(tx pgx.Tx, login string, password string, name string) (pgconn.CommandTag, error) {
	exec, err := tx.Exec(
		context.Background(),
		"INSERT INTO users (login, password, name) VALUES ($1, $2, $3)",
		login,
		password,
		name)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, internalErrors.ErrUniqueViolation
		}
		return nil, err
	}
	return exec, nil
}

func UpdateToken(dbr *Repository, login string, token string, expiredAt int64) (pgconn.CommandTag, error) {
	exec, err := dbr.ExecContext(
		context.Background(),
		" UPDATE users SET access_token = $1, access_token_expired_at = $2 WHERE login = $3",
		token,
		expiredAt,
		login)

	if err != nil {
		return nil, internalErrors.NewDBError(err)
	}
	return exec, nil
}

func AddOrder(dbr *Repository, number string, accrual float32, statusID int, uploadedAt string, userID int) (pgconn.CommandTag, error) {
	exec, err := dbr.ExecContext(
		context.Background(),
		`INSERT INTO "order" (number, accrual, status_id, uploaded_at, user_id) VALUES ($1, $2, $3, $4, $5)`,
		number,
		accrual,
		statusID,
		uploadedAt,
		userID)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, internalErrors.ErrUniqueViolation
		}
		return nil, err
	}
	return exec, nil
}

func AddWithdraw(tx pgx.Tx, orderNumber string, sum float32, processedAt string, login string) (pgconn.CommandTag, error) {
	exec, err := tx.Exec(
		context.Background(),
		`INSERT INTO withdraw (sum, "order", processed_at, user_id) VALUES ($1, $2, $3, (SELECT id FROM users WHERE login = $4))`,
		sum,
		orderNumber,
		processedAt,
		login)
	return exec, err
}

func OperatingBalance(tx pgx.Tx, sum float32, operator string, login string) (pgconn.CommandTag, error) {
	exec, err := tx.Exec(
		context.Background(),
		`
			UPDATE balance 
			SET "value" = "value" `+operator+` $1 
			WHERE user_id = (
				SELECT id 
				FROM users 
				WHERE login = $2 LIMIT 1
			);`,
		sum,
		login)

	if err != nil {
		return nil, err
	}
	return exec, nil
}

func UpdateOrderAccrual(tx pgx.Tx, number string, accrual float32, statusID int) (pgconn.CommandTag, error) {
	exec, err := tx.Exec(
		context.Background(),
		`
			UPDATE "order" 
			SET "accrual" = $1, "status_id" = $2 
			WHERE number = $3`,
		accrual,
		statusID,
		number)

	if err != nil {
		return nil, err
	}
	return exec, nil
}

func UpdateOrderCheckAccrualAfter(dbr *Repository, number string, checkAccrualAfter int64) (pgconn.CommandTag, error) {
	exec, err := dbr.ExecContext(
		context.Background(),
		`
			UPDATE "order" 
			SET "check_accrual_after" = $1
			WHERE number = $2`,
		checkAccrualAfter,
		number)

	if err != nil {
		return nil, err
	}
	return exec, nil
}

func GetOrder(tx pgx.Tx, number string) (user.Order, error) {
	var o user.Order
	row := tx.QueryRow(
		context.Background(),
		`SELECT o.id, number, accrual::DOUBLE PRECISION, s.alias, uploaded_at, u.login
				FROM "order" o
				LEFT JOIN users u ON u.id = o.user_id
				LEFT JOIN status s ON s.id = o.status_id
				WHERE number = $1`,
		number)

	err := row.Scan(&o.ID, &o.Number, &o.Accrual, &o.Status, &o.UploadedAt, &o.UserLogin)
	return o, err
}

func GetBalance(tx pgx.Tx, login string) (balance.Balance, error) {
	var b balance.Balance
	row := tx.QueryRow(
		context.Background(),
		`SELECT b.id,
					   "value"::DOUBLE PRECISION as balance,
					   COALESCE(sum(w.sum) OVER (PARTITION BY b.id), 0::DOUBLE PRECISION) as sum
				FROM balance b
						 LEFT JOIN users u ON u.id = b.user_id
						 LEFT JOIN withdraw w ON w.user_id = b.user_id
				WHERE u.login = $1
				LIMIT 1`,
		login)
	err := row.Scan(&b.ID, &b.Current, &b.Withdrawn)
	return b, err
}

func AddBalance(tx pgx.Tx, login string) (pgconn.CommandTag, error) {
	exec, err := tx.Exec(
		context.Background(),
		"INSERT INTO balance (user_id) (SELECT id FROM users WHERE login = $1)",
		login)
	return exec, err
}
