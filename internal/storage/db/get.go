package db

import (
	"context"
	"time"

	"github.com/lib/pq"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
	"github.com/ramil063/firstgodiplom/internal/logger"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository"
)

// GetUser получить пользователя
func (s *Storage) GetUser(ctx context.Context, login string) (user.User, error) {
	var u user.User
	row := repository.DBRepository.QueryRowContext(
		ctx,
		"SELECT id, login, password, name FROM users WHERE login = $1",
		login)
	err := row.Scan(&u.ID, &u.Login, &u.PasswordHash, &u.Name)
	return u, err
}

// GetAccessTokenData получить данные токена авторизации
func (s *Storage) GetAccessTokenData(ctx context.Context, token string) (user.AccessTokenData, error) {
	var t user.AccessTokenData
	query := "SELECT login, access_token, access_token_expired_at FROM users WHERE access_token = $1"
	row := repository.DBRepository.QueryRowContext(
		ctx,
		query,
		token)
	err := row.Scan(&t.Login, &t.AccessToken, &t.AccessTokenExpiredAt)
	return t, err
}

// GetOrder получить заказ
func (s *Storage) GetOrder(ctx context.Context, number string) (user.Order, error) {
	var o user.Order
	row := repository.DBRepository.QueryRowContext(
		ctx,
		`SELECT o.id, number, accrual::DECIMAL, s.alias, uploaded_at, u.login
				FROM "order" o
				LEFT JOIN users u ON u.id = o.user_id
				LEFT JOIN status s ON s.id = o.status_id
				WHERE number = $1`,
		number)
	err := row.Scan(&o.ID, &o.Number, &o.Accrual, &o.Status, &o.UploadedAt, &o.UserLogin)
	return o, err
}

// GetOrders получить заказы
func (s *Storage) GetOrders(ctx context.Context, login string) ([]user.Order, error) {
	var res []user.Order
	rows, err := repository.DBRepository.QueryContext(
		ctx,
		`SELECT number, accrual::DECIMAL, s.alias, uploaded_at
				FROM "order" o
				LEFT JOIN users u ON u.id = o.user_id
				LEFT JOIN status s ON s.id = o.status_id
				WHERE u.login = $1`,
		login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accrual float32
	var number, status, uploadedAt string
	for rows.Next() {
		err = rows.Scan(&number, &accrual, &status, &uploadedAt)
		if err != nil {
			logger.WriteErrorLog(err.Error())
			return nil, err
		}
		res = append(res, user.Order{
			Number:     number,
			Status:     status,
			Accrual:    accrual,
			UploadedAt: uploadedAt,
		})
	}

	err = rows.Err()
	if err != nil {
		logger.WriteErrorLog(err.Error())
	}
	return res, err
}

// GetAllOrdersInStatuses получить все заказы в статусах
func (s *Storage) GetAllOrdersInStatuses(ctx context.Context, statuses []int) ([]user.OrderCheckAccrual, error) {
	var res []user.OrderCheckAccrual

	rows, err := repository.DBRepository.QueryContext(
		ctx,
		`SELECT number, accrual::DECIMAL, s.alias
				FROM "order" o
				LEFT JOIN status s ON s.id = o.status_id
				WHERE status_id = ANY($1)
				AND (check_accrual_after IS NULL OR check_accrual_after <= $2)
				LIMIT 50`,
		pq.Array(statuses),
		time.Now().Unix())

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accrual float32
	var number, status string
	for rows.Next() {
		err = rows.Scan(&number, &accrual, &status)
		if err != nil {
			logger.WriteErrorLog(err.Error())
			return nil, err
		}
		res = append(res, user.OrderCheckAccrual{
			Number:  number,
			Accrual: accrual,
			Status:  status,
		})
	}

	err = rows.Err()
	if err != nil {
		logger.WriteErrorLog(err.Error())
	}
	return res, err
}

// GetBalance получить баланс пользователя
func (s *Storage) GetBalance(ctx context.Context, login string) (balance.Balance, error) {
	var b balance.Balance
	row := repository.DBRepository.QueryRowContext(
		ctx,
		`SELECT
					"value"::DECIMAL AS balance,
					COALESCE(SUM(w.sum) OVER(PARTITION BY b.id), 0::DECIMAL) AS sum
				FROM balance b
						 INNER JOIN users u ON u.id = b.user_id
						 LEFT JOIN withdraw w ON w.user_id = b.user_id
				WHERE u.login = $1
				LIMIT 1`,
		login)
	err := row.Scan(&b.Current, &b.Withdrawn)
	return b, err
}

// GetWithdrawals получение всех списаний баллов
func (s *Storage) GetWithdrawals(ctx context.Context, login string) ([]balance.Withdraw, error) {
	var res []balance.Withdraw
	rows, err := repository.DBRepository.QueryContext(
		ctx,
		`SELECT "order", "sum", processed_at
				FROM withdraw w
				INNER JOIN users u ON u.id = w.user_id
				WHERE u.login = $1
				ORDER BY processed_at DESC`,
		login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sum float32
	var number, processedAt string
	for rows.Next() {
		err = rows.Scan(&number, &sum, &processedAt)
		if err != nil {
			logger.WriteErrorLog(err.Error())
			return nil, err
		}
		res = append(res, balance.Withdraw{
			OrderNumber: number,
			Sum:         sum,
			ProcessedAt: processedAt,
		})
	}

	err = rows.Err()
	if err != nil {
		logger.WriteErrorLog(err.Error())
	}
	return res, err
}
