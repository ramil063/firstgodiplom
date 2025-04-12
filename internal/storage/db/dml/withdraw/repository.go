package withdraw

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// AddWithdraw добавить списание средств
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
