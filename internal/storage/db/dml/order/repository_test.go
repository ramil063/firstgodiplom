package order

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository"
	repository2 "github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository/mocks"
)

func TestAddOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		number     string
		accrual    float32
		statusID   int
		uploadedAt string
		userID     int
	}{
		{
			name:       "test 1",
			number:     "1",
			accrual:    100.0,
			statusID:   1,
			uploadedAt: "2025-04-15",
			userID:     1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			dbr := repository.Repository{Pool: poolMock}

			expectedCommandTag := pgconn.CommandTag("INSERT 0 1")
			poolMock.EXPECT().
				Exec(
					context.Background(),
					`INSERT INTO "order" (number, accrual, status_id, uploaded_at, user_id) VALUES ($1, $2, $3, $4, $5)`,
					tt.number,
					tt.accrual,
					tt.statusID,
					tt.uploadedAt,
					tt.userID).
				Return(expectedCommandTag, nil)

			_, err := AddOrder(context.Background(), &dbr, tt.number, tt.accrual, tt.statusID, tt.uploadedAt, tt.userID)
			assert.NoError(t, err)
		})
	}
}

func TestGetOrder(t *testing.T) {
	tests := []struct {
		name   string
		number string
		order  user.Order
	}{
		{
			name:   "test 1",
			number: "1",
			order: user.Order{
				ID:         1,
				Number:     "1",
				Status:     "NEW",
				Accrual:    100.1,
				UploadedAt: "2025-04-14",
				UserLogin:  "ramil",
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
			expectedOrder := tt.order

			rows := mock.NewRows([]string{"id", "number", "accrual", "alias", "uploaded_at", "login"}).
				AddRow(
					expectedOrder.ID,
					expectedOrder.Number,
					expectedOrder.Accrual,
					expectedOrder.Status,
					expectedOrder.UploadedAt,
					expectedOrder.UserLogin)

			mock.ExpectQuery(`^SELECT o.id, number, accrual::DECIMAL, s.alias, uploaded_at, u.login*
				FROM "order" o*
				LEFT JOIN users u ON u.id = o.user_id*
				LEFT JOIN status s ON s.id = o.status_id*
				WHERE number = \$1`).
				WithArgs(tt.number).
				WillReturnRows(rows)

			mock.ExpectCommit()

			tx, err := mock.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			got, err := GetOrder(context.Background(), tx, tt.number)
			assert.Equal(t, expectedOrder.ID, got.ID)
			assert.Equal(t, expectedOrder.Number, got.Number)
			assert.Equal(t, expectedOrder.Status, got.Status)
			assert.Equal(t, expectedOrder.Accrual, got.Accrual)
			assert.Equal(t, expectedOrder.UploadedAt, got.UploadedAt)
			assert.Equal(t, expectedOrder.UserLogin, got.UserLogin)
			assert.NoError(t, err)
		})
	}
}

func TestUpdateOrderAccrual(t *testing.T) {

	tests := []struct {
		name     string
		number   string
		accrual  float32
		statusID int
	}{
		{
			name:     "test 1",
			number:   "1",
			accrual:  100.0,
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

			mock.ExpectBegin()
			expectedCommandTag := pgconn.CommandTag("UPDATE 0 1")

			mock.ExpectExec(`UPDATE "order".*SET "accrual" = \$1, "status_id" = \$2.*WHERE number = \$3.*`).
				WithArgs(
					tt.accrual,
					tt.statusID,
					tt.number).
				WillReturnResult(expectedCommandTag)
			mock.ExpectCommit()

			tx, err := mock.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}
			_, err = UpdateOrderAccrual(context.Background(), tx, tt.number, tt.accrual, tt.statusID)
			assert.NoError(t, err)
		})
	}
}

func TestUpdateOrderCheckAccrualAfter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name              string
		number            string
		checkAccrualAfter int64
	}{
		{
			name:              "test 1",
			number:            "1",
			checkAccrualAfter: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			dbr := repository.Repository{Pool: poolMock}

			expectedCommandTag := pgconn.CommandTag("UPDATE 0 1")
			poolMock.EXPECT().
				Exec(
					context.Background(),
					`UPDATE "order" SET "check_accrual_after" = $1 WHERE number = $2`,
					tt.checkAccrualAfter,
					tt.number).
				Return(expectedCommandTag, nil)

			_, err := UpdateOrderCheckAccrualAfter(context.Background(), &dbr, tt.number, tt.checkAccrualAfter)
			assert.NoError(t, err)
		})
	}
}
