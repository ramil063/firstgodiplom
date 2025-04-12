package order

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	internalErrors "github.com/ramil063/firstgodiplom/internal/errors"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository"
)

// AddOrder добавление заказа
func AddOrder(dbr *repository.Repository, number string, accrual float32, statusID int, uploadedAt string, userID int) (pgconn.CommandTag, error) {
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

// UpdateOrderAccrual обновление начисления в заказе
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

// UpdateOrderCheckAccrualAfter обновление поля даты следующей проверки заказа в акруал
func UpdateOrderCheckAccrualAfter(dbr *repository.Repository, number string, checkAccrualAfter int64) (pgconn.CommandTag, error) {
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

// GetOrder получение заказа
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
