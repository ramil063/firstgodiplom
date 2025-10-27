package db

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/agent/accrual/storage"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository"
	repository2 "github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository/mocks"
)

func TestStorage_UpdateOrderAccrual(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name             string
		orderFromAccrual storage.Order
		expectedOrder    user.Order
		expectedBalance  balance.Balance
		statusID         int
	}{
		{
			name: "test 1",
			orderFromAccrual: storage.Order{
				Status:    "REGISTERED",
				Order:     "1",
				Accrual:   100.1,
				UserLogin: "ramil",
			},
			expectedOrder: user.Order{
				ID:         1,
				Number:     "1",
				Status:     "NEW",
				Accrual:    100.1,
				UploadedAt: "2025-04-14",
				UserLogin:  "ramil",
			},
			expectedBalance: balance.Balance{
				ID:        1,
				Current:   100.1,
				Withdrawn: 100.1,
			},
			statusID: 1,
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
			expectedCommandTag := pgconn.CommandTag("UPDATE 0 1")
			mock.ExpectExec(`.*UPDATE "order",*`).
				WithArgs(
					tt.orderFromAccrual.Accrual,
					tt.statusID,
					tt.orderFromAccrual.Order).
				WillReturnResult(expectedCommandTag)

			rows := mock.NewRows([]string{"id", "number", "accrual", "alias", "uploaded_at", "login"}).
				AddRow(
					tt.expectedOrder.ID,
					tt.expectedOrder.Number,
					tt.expectedOrder.Accrual,
					tt.expectedOrder.Status,
					tt.expectedOrder.UploadedAt,
					tt.expectedOrder.UserLogin)
			mock.ExpectQuery(`.*SELECT o.id, number, accrual::DECIMAL`).
				WithArgs(tt.orderFromAccrual.Order).
				WillReturnRows(rows)

			rowsBalance := mock.NewRows([]string{"id", "current", "withdrawn"}).
				AddRow(
					tt.expectedBalance.ID,
					tt.expectedBalance.Current,
					tt.expectedBalance.Withdrawn)
			mock.ExpectQuery(`.*FROM balance.*`).
				WithArgs(tt.orderFromAccrual.UserLogin).
				WillReturnRows(rowsBalance)

			rows = mock.NewRows([]string{"value"}).
				AddRow(tt.orderFromAccrual.Accrual)
			mock.ExpectQuery(`.*UPDATE balance.*`).
				WithArgs(
					tt.orderFromAccrual.Accrual,
					tt.orderFromAccrual.UserLogin).
				WillReturnRows(rows)

			mock.ExpectCommit()
			s := &Storage{}
			err = s.UpdateOrderAccrual(context.Background(), tt.orderFromAccrual)
			assert.NoError(t, err)
		})
	}
}

func TestStorage_UpdateOrderCheckAccrualAfter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		orderNumber string
	}{
		{
			name:        "test 1",
			orderNumber: "1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			repository.DBRepository = repository.Repository{Pool: poolMock}

			expectedCommandTag := pgconn.CommandTag("UPDATE 0 1")
			poolMock.EXPECT().
				Exec(
					context.Background(),
					`UPDATE "order" SET "check_accrual_after" = $1 WHERE number = $2`,
					time.Now().Unix(),
					tt.orderNumber).
				Return(expectedCommandTag, nil)
			s := &Storage{}
			err := s.UpdateOrderCheckAccrualAfter(context.Background(), tt.orderNumber)
			assert.NoError(t, err)
		})
	}
}

func TestStorage_UpdateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name      string
		login     string
		token     auth.Token
		expiredAt int64
	}{
		{
			name:      "test 1",
			login:     "ramil",
			token:     auth.Token{Token: "token"},
			expiredAt: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			repository.DBRepository = repository.Repository{Pool: poolMock}

			expectedCommandTag := pgconn.CommandTag("UPDATE 0 1")
			poolMock.EXPECT().
				Exec(
					context.Background(),
					`UPDATE users SET access_token = $1, access_token_expired_at = $2 WHERE login = $3`,
					tt.token.Token,
					tt.expiredAt,
					tt.login).
				Return(expectedCommandTag, nil)
			s := &Storage{}
			err := s.UpdateToken(context.Background(), tt.login, tt.token, tt.expiredAt)
			assert.NoError(t, err)
		})
	}
}
