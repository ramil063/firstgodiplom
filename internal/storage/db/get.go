package db

import (
	"context"
	"time"

	"github.com/lib/pq"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
	"github.com/ramil063/firstgodiplom/internal/logger"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml"
)

func (s *Storage) GetUser(login string) (user.User, error) {
	var u user.User
	row := dml.DBRepository.QueryRowContext(
		context.Background(),
		"SELECT id, login, password, name FROM users WHERE login = $1",
		login)
	err := row.Scan(&u.ID, &u.Login, &u.PasswordHash, &u.Name)
	return u, err
}

func (s *Storage) GetAccessTokenData(token string) (user.AccessTokenData, error) {
	var t user.AccessTokenData
	query := "SELECT login, access_token, access_token_expired_at FROM users WHERE access_token = $1"
	row := dml.DBRepository.QueryRowContext(
		context.Background(),
		query,
		token)
	err := row.Scan(&t.Login, &t.AccessToken, &t.AccessTokenExpiredAt)
	return t, err
}

func (s *Storage) GetOrder(number string) (user.Order, error) {
	var o user.Order
	row := dml.DBRepository.QueryRowContext(
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

func (s *Storage) GetOrders(login string) ([]user.Order, error) {
	var res []user.Order
	rows, err := dml.DBRepository.QueryContext(
		context.Background(),
		`SELECT number, accrual::DOUBLE PRECISION, s.alias, uploaded_at
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
			continue
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

func (s *Storage) GetAllOrdersInStatuses(statuses []int) ([]user.OrderCheckAccrual, error) {
	var res []user.OrderCheckAccrual

	rows, err := dml.DBRepository.QueryContext(
		context.Background(),
		`SELECT number, accrual::DOUBLE PRECISION, s.alias
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
			continue
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

func (s *Storage) GetBalance(login string) (balance.Balance, error) {
	var b balance.Balance
	row := dml.DBRepository.QueryRowContext(
		context.Background(),
		`SELECT max("value") as balance, sum(w.sum) as sum
				FROM balance b
						 LEFT JOIN users u ON u.id = b.user_id
						 LEFT JOIN public."order" o on u.id = o.user_id
						 LEFT JOIN withdraw w ON w.order_id = o.id
				WHERE u.login = $1
				GROUP BY b."user_id"`,
		login)
	err := row.Scan(&b.Current, &b.Withdrawn)
	return b, err
}

func (s *Storage) GetCounter(name string) (int64, error) {
	row := dml.DBRepository.QueryRowContext(context.Background(), "SELECT value FROM counter WHERE name = $1", name)
	var selectedValue int64
	err := row.Scan(&selectedValue)

	return selectedValue, err
}

func (s *Storage) GetWithdrawals(login string) ([]balance.Withdraw, error) {
	var res []balance.Withdraw
	rows, err := dml.DBRepository.QueryContext(
		context.Background(),
		`SELECT o.number, "sum", processed_at
				FROM withdraw w
				LEFT JOIN "order" o ON w.order_id = o.id
				LEFT JOIN users u ON u.id = o.user_id
				WHERE u.login = $1
				ORDER BY processed_at DESC`,
		login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sum int
	var number, processedAt string
	for rows.Next() {
		err = rows.Scan(&number, &sum, &processedAt)
		if err != nil {
			logger.WriteErrorLog(err.Error())
			continue
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
