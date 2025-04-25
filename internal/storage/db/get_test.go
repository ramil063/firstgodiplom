package db

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository"
	repository2 "github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository/mocks"
)

type mockRow struct {
	values []interface{}
	err    error
}

func (m *mockRow) Scan(dest ...interface{}) error {
	if m.err != nil {
		return m.err
	}
	for i := range dest {
		if i >= len(m.values) {
			return errors.New("not enough values to scan")
		}
		val := reflect.ValueOf(dest[i]).Elem()
		val.Set(reflect.ValueOf(m.values[i]))
	}
	return nil
}

func TestStorage_GetAccessTokenData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name              string
		query             string
		token             string
		expectedTokenData user.AccessTokenData
	}{
		{
			name:  "test 1",
			query: `SELECT login, access_token, access_token_expired_at FROM users WHERE access_token = $1`,
			token: "123",
			expectedTokenData: user.AccessTokenData{
				Login:                "ramil",
				AccessToken:          "123",
				AccessTokenExpiredAt: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()
			repository.DBRepository = repository.Repository{Pool: poolMock}

			poolMock.EXPECT().
				QueryRow(
					context.Background(),
					tt.query,
					tt.token).
				Return(&mockRow{
					values: []interface{}{
						tt.expectedTokenData.Login,
						tt.expectedTokenData.AccessToken,
						tt.expectedTokenData.AccessTokenExpiredAt,
					},
				})

			s := &Storage{}
			_, err = s.GetAccessTokenData(context.Background(), tt.token)
			assert.NoError(t, err)
		})
	}
}

func TestStorage_GetAllOrdersInStatuses(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		want     []user.OrderCheckAccrual
		statuses []int
	}{
		{
			name: "test 1",
			want: []user.OrderCheckAccrual{{
				Number:  "1",
				Accrual: 100.0,
				Status:  "NEW",
			}},
			statuses: []int{1, 2},
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

			rows := mock.NewRows([]string{"number", "accrual", "status"}).
				AddRow(tt.want[0].Number, tt.want[0].Accrual, tt.want[0].Status)

			mock.ExpectQuery(`.*SELECT number\, accrual\:\:DECIMAL\, s.alias.*`).
				WithArgs(
					pq.Array(tt.statuses),
					time.Now().Unix()).
				WillReturnRows(rows)

			s := &Storage{}

			_, err = s.GetAllOrdersInStatuses(context.Background(), tt.statuses)
			assert.NoError(t, err)
		})
	}
}

func TestStorage_GetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name  string
		login string
		want  balance.Balance
	}{
		{
			name:  "test 1",
			login: "ramil",
			want: balance.Balance{
				Current:   100.1,
				Withdrawn: 100.1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()
			repository.DBRepository = repository.Repository{Pool: poolMock}

			poolMock.EXPECT().
				QueryRow(
					context.Background(),
					gomock.Any(),
					tt.login).
				Return(&mockRow{
					values: []interface{}{
						tt.want.Current,
						tt.want.Withdrawn,
					},
				})
			s := &Storage{}
			_, err = s.GetBalance(context.Background(), tt.login)
			assert.NoError(t, err)
		})
	}
}

func TestStorage_GetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name   string
		number string
		want   user.Order
	}{
		{
			name:   "test 1",
			number: "1",
			want: user.Order{
				ID:         1,
				Number:     "1",
				Accrual:    100.1,
				Status:     "NEW",
				UploadedAt: "2025-04-17",
				UserLogin:  "ramil",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()
			repository.DBRepository = repository.Repository{Pool: poolMock}

			poolMock.EXPECT().
				QueryRow(
					context.Background(),
					gomock.Any(),
					tt.number).
				Return(&mockRow{
					values: []interface{}{
						tt.want.ID,
						tt.want.Number,
						tt.want.Accrual,
						tt.want.Status,
						tt.want.UploadedAt,
						tt.want.UserLogin,
					},
				})
			s := &Storage{}
			_, err = s.GetOrder(context.Background(), tt.number)
			assert.NoError(t, err)
		})
	}
}

func TestStorage_GetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name  string
		login string
		want  []user.Order
	}{
		{
			name:  "test 1",
			login: "ramil",
			want: []user.Order{{
				Number:     "1",
				Accrual:    100.1,
				Status:     "NEW",
				UploadedAt: "2025-04-17",
			}},
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

			rows := mock.NewRows([]string{"number", "accrual", "status", "status"}).
				AddRow(tt.want[0].Number, tt.want[0].Accrual, tt.want[0].Status, tt.want[0].UploadedAt)

			mock.ExpectQuery(`.*SELECT number\, accrual\:\:DECIMAL\, s.alias.*`).
				WithArgs(tt.login).
				WillReturnRows(rows)

			s := &Storage{}
			_, err = s.GetOrders(context.Background(), tt.login)
			assert.NoError(t, err)
		})
	}
}

func TestStorage_GetUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name  string
		login string
		want  user.User
	}{
		{
			name:  "test 1",
			login: "ramil",
			want: user.User{
				ID:           1,
				Login:        "ramil",
				PasswordHash: "123",
				Name:         "ramil",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			mock, err := pgxmock.NewPool()
			if err != nil {
				t.Fatal(err)
			}
			defer mock.Close()
			repository.DBRepository = repository.Repository{Pool: poolMock}

			poolMock.EXPECT().
				QueryRow(
					context.Background(),
					gomock.Any(),
					tt.login).
				Return(&mockRow{
					values: []interface{}{
						tt.want.ID,
						tt.want.Login,
						tt.want.PasswordHash,
						tt.want.Name,
					},
				})
			s := &Storage{}
			_, err = s.GetUser(context.Background(), tt.login)
			assert.NoError(t, err)
		})
	}
}

func TestStorage_GetWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name  string
		login string
		want  []balance.Withdraw
	}{
		{
			name:  "test 1",
			login: "ramil",
			want: []balance.Withdraw{{
				OrderNumber: "1",
				Sum:         100.1,
				ProcessedAt: "2025-04-17",
			}},
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

			rows := mock.NewRows([]string{"number", "sum", "processed_at"}).
				AddRow(tt.want[0].OrderNumber, tt.want[0].Sum, tt.want[0].ProcessedAt)

			mock.ExpectQuery(`SELECT "order"\, "sum"\, processed_at.*`).
				WithArgs(tt.login).
				WillReturnRows(rows)

			s := &Storage{}
			_, err = s.GetWithdrawals(context.Background(), tt.login)
			assert.NoError(t, err)
		})
	}
}
