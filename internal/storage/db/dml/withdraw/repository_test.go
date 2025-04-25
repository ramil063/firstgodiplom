package withdraw

import (
	"context"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
)

func TestAddWithdraw(t *testing.T) {
	tests := []struct {
		name        string
		orderNumber string
		sum         float32
		processedAt string
		login       string
	}{
		{
			name:        "test 1",
			orderNumber: "1",
			sum:         100,
			processedAt: "2025-04-14",
			login:       "ramil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()

			mock.ExpectBegin()
			expectedCommandTag := pgconn.CommandTag("INSERT 0 1")
			mock.ExpectExec(`^INSERT INTO withdraw \(sum, "order", processed_at, user_id\) VALUES \(\$1, \$2, \$3, \(SELECT id FROM users WHERE login = \$4\)\)`).
				WithArgs(
					tt.sum,
					tt.orderNumber,
					tt.processedAt,
					tt.login).
				WillReturnResult(expectedCommandTag)
			mock.ExpectCommit()

			tx, err := mock.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			_, err = AddWithdraw(context.Background(), tx, tt.orderNumber, tt.sum, tt.processedAt, tt.login)
			assert.NoError(t, err)
		})
	}
}
