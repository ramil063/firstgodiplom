package balance

import (
	"context"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
)

func TestAddBalance(t *testing.T) {

	tests := []struct {
		name  string
		login string
	}{
		{
			name:  "test 1",
			login: "ramil",
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
			mock.ExpectExec(`^INSERT INTO balance \(user_id\) \(SELECT id FROM users WHERE login = \$1\)$`).
				WithArgs(tt.login).
				WillReturnResult(expectedCommandTag)
			mock.ExpectCommit()

			tx, err := mock.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			_, err = AddBalance(tx, tt.login)
			assert.NoError(t, err)
		})
	}
}

func TestGetBalance(t *testing.T) {

	tests := []struct {
		name            string
		login           string
		expectedBalance balance.Balance
	}{
		{
			name:  "test 1",
			login: "ramil",
			expectedBalance: balance.Balance{
				ID:        1,
				Current:   100.50,
				Withdrawn: 30.25,
			},
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
			expectedBalance := tt.expectedBalance

			rows := mock.NewRows([]string{"id", "balance", "sum"}).
				AddRow(expectedBalance.ID, expectedBalance.Current, expectedBalance.Withdrawn)

			mock.ExpectQuery(`SELECT b.id, "value"::DECIMAL as balance, .*WHERE u.login = \$1`).
				WithArgs(tt.login).
				WillReturnRows(rows)

			mock.ExpectCommit()

			tx, err := mock.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			got, err := GetBalanceForUpdate(tx, tt.login)
			assert.Equal(t, got.ID, expectedBalance.ID)
			assert.Equal(t, got.Current, expectedBalance.Current)
			assert.Equal(t, got.Withdrawn, expectedBalance.Withdrawn)
			assert.NoError(t, err)
		})
	}
}

func TestOperatingBalance(t *testing.T) {

	tests := []struct {
		name     string
		login    string
		operator string
		sum      float32
	}{
		{
			name:     "test 1",
			login:    "ramil",
			operator: "+",
			sum:      100,
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

			rows := mock.NewRows([]string{"value"}).
				AddRow(tt.sum)
			mock.ExpectQuery(`UPDATE balance.*SET "value" = "value" \+ \$1.*WHERE user_id = \(.*SELECT id.*FROM users.*WHERE login = \$2.*\)*RETURNING "value";`).
				WithArgs(tt.sum, tt.login).
				WillReturnRows(rows)
			mock.ExpectCommit()

			tx, err := mock.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			_, err = OperatingBalance(tx, tt.sum, tt.operator, tt.login)
			assert.NoError(t, err)
		})
	}
}
