package db

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository"
	repository2 "github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository/mocks"
)

func TestStorage_AddOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name         string
		number       string
		statusID     int
		tokenData    user.AccessTokenData
		userID       int
		accrual      float32
		passwordHash string
	}{
		{
			name:     "test 1",
			number:   "111",
			statusID: 1,
			tokenData: user.AccessTokenData{
				Login:                "ramil",
				AccessToken:          "123123",
				AccessTokenExpiredAt: 1,
			},
			userID:       1,
			accrual:      0,
			passwordHash: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			repository.DBRepository = repository.Repository{Pool: poolMock}

			poolMock.EXPECT().
				QueryRow(
					context.Background(),
					gomock.Any(),
					tt.tokenData.Login).
				Return(&mockRow{
					values: []interface{}{
						tt.userID,
						tt.tokenData.Login,
						tt.passwordHash,
						tt.tokenData.Login,
					},
				})
			expectedCommandTag := pgconn.CommandTag("INSERT 0 1")
			poolMock.EXPECT().
				Exec(
					context.Background(),
					`INSERT INTO "order" (number, accrual, status_id, uploaded_at, user_id) VALUES ($1, $2, $3, $4, $5)`,
					tt.number,
					tt.accrual,
					tt.statusID,
					gomock.Any(),
					tt.userID).
				Return(expectedCommandTag, nil)

			s := &Storage{}
			err := s.AddOrder(tt.number, tt.tokenData)
			assert.NoError(t, err)
		})
	}
}

func TestStorage_AddUserData(t *testing.T) {
	tests := []struct {
		name         string
		register     user.Register
		passwordHash string
	}{
		{
			name: "test 1",
			register: user.Register{
				Login:    "ramil",
				Password: "123123",
				Name:     "ramil",
			},
			passwordHash: "123123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()
			repository.DBRepository = repository.Repository{Pool: mock}

			mock.ExpectBegin()
			expectedCommandTag := pgconn.CommandTag("INSERT 0 1")
			mock.ExpectExec(`^INSERT INTO users \(login, password, name\) VALUES \(\$1, \$2, \$3\)$`).
				WithArgs(
					tt.register.Login,
					tt.register.Password,
					tt.register.Name).
				WillReturnResult(expectedCommandTag)

			expectedCommandTag = pgconn.CommandTag("INSERT 0 1")
			mock.ExpectExec(`^INSERT INTO balance \(user_id\) \(SELECT id FROM users WHERE login = \$1\)$`).
				WithArgs(tt.register.Login).
				WillReturnResult(expectedCommandTag)
			mock.ExpectCommit()

			s := &Storage{}
			err = s.AddUserData(tt.register, tt.passwordHash)
			assert.NoError(t, err)
		})
	}
}

func TestStorage_AddWithdrawFromBalance(t *testing.T) {

	tests := []struct {
		name            string
		login           string
		withdraw        balance.Withdraw
		expectedBalance balance.Balance
	}{
		{
			name:  "test 1",
			login: "ramil",
			withdraw: balance.Withdraw{
				OrderNumber: "1",
				Sum:         1,
				ProcessedAt: "2025-04-17",
				UserLogin:   "ramil",
			},
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
			repository.DBRepository = repository.Repository{Pool: mock}

			mock.ExpectBegin()
			rows := mock.NewRows([]string{"id", "balance", "sum"}).
				AddRow(tt.expectedBalance.ID, tt.expectedBalance.Current, tt.expectedBalance.Withdrawn)

			mock.ExpectQuery(`SELECT b.id, "value"::DECIMAL as balance, COALESCE.*WHERE u.login = \$1`).
				WithArgs(tt.login).
				WillReturnRows(rows)

			expectedCommandTag := pgconn.CommandTag("INSERT 0 1")
			mock.ExpectExec(`^INSERT INTO withdraw \(sum, "order", processed_at, user_id\) VALUES \(\$1, \$2, \$3, \(SELECT id FROM users WHERE login = \$4\)\)`).
				WithArgs(
					tt.withdraw.Sum,
					tt.withdraw.OrderNumber,
					pgxmock.AnyArg(),
					tt.withdraw.UserLogin,
				).
				WillReturnResult(expectedCommandTag)

			rows = mock.NewRows([]string{"value"}).
				AddRow(tt.withdraw.Sum)
			mock.ExpectQuery(`UPDATE balance.*SET "value" = "value" \- \$1.*WHERE user_id = \(.*SELECT id.*FROM users.*WHERE login = \$2.*\)*RETURNING "value";`).
				WithArgs(tt.withdraw.Sum, tt.login).
				WillReturnRows(rows)
			mock.ExpectCommit()

			s := &Storage{}
			err = s.AddWithdrawFromBalance(tt.withdraw, tt.login)
			assert.NoError(t, err)
		})
	}
}
