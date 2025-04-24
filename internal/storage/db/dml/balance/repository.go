package balance

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
)

// OperatingBalance обновление баланса
func OperatingBalance(tx pgx.Tx, sum float32, operator string, login string) (float32, error) {
	row := tx.QueryRow(
		context.Background(),
		`
			UPDATE balance 
			SET "value" = "value" `+operator+` $1 
			WHERE user_id = (
				SELECT id 
				FROM users 
				WHERE login = $2 LIMIT 1
			)
			RETURNING "value";`,
		sum,
		login)

	var newBalance float32
	err := row.Scan(&newBalance)
	return newBalance, err
}

// GetBalanceForUpdate получение баланса
func GetBalanceForUpdate(tx pgx.Tx, login string) (balance.Balance, error) {
	var b balance.Balance
	row := tx.QueryRow(
		context.Background(),
		`SELECT b.id,
					   "value"::DECIMAL as balance,
					   (SELECT COALESCE(SUM(w.sum), 0::DECIMAL) 
						FROM withdraw w 
						WHERE w.user_id = b.user_id) as sum
				FROM balance b
						 LEFT JOIN users u ON u.id = b.user_id
						 LEFT JOIN withdraw w ON w.user_id = b.user_id
				WHERE u.login = $1
				LIMIT 1
				FOR UPDATE OF b`,
		login)
	err := row.Scan(&b.ID, &b.Current, &b.Withdrawn)
	return b, err
}

// AddBalance добавление баланса
func AddBalance(tx pgx.Tx, login string) (pgconn.CommandTag, error) {
	exec, err := tx.Exec(
		context.Background(),
		"INSERT INTO balance (user_id) (SELECT id FROM users WHERE login = $1)",
		login)
	return exec, err
}
