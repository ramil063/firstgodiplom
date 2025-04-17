package repository

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"

	repository2 "github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository/mocks"
)

func TestRepository_ExecContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name            string
		ctx             context.Context
		query           string
		commandTagQuery string
		token           string
		expiredAt       int64
		login           string
	}{
		{
			name:            "success",
			ctx:             context.Background(),
			query:           "UPDATE users SET access_token = $1, access_token_expired_at = $2 WHERE login = $3",
			commandTagQuery: "UPDATE 0 1",
			token:           "token",
			expiredAt:       10,
			login:           "ramil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			dbr := Repository{Pool: poolMock}

			expectedCommandTag := pgconn.CommandTag(tt.commandTagQuery)
			poolMock.EXPECT().
				Exec(
					context.Background(),
					tt.query,
					tt.token,
					tt.expiredAt,
					tt.login).
				Return(expectedCommandTag, nil)

			_, err := dbr.ExecContext(tt.ctx, tt.query, tt.token, tt.expiredAt, tt.login)
			assert.NoError(t, err)
		})
	}
}

func TestRepository_PingContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
	}{
		{"test 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			dbr := Repository{Pool: poolMock}

			poolMock.EXPECT().
				Ping(context.Background()).
				Return(nil)
			err := dbr.PingContext(context.Background())
			assert.NoError(t, err)
		})
	}
}

func TestRepository_QueryContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name   string
		number string
		query  string
	}{
		{
			name:   "test 1",
			number: "1",
			query: `^SELECT o.id, number, accrual::DOUBLE PRECISION, s.alias, uploaded_at, u.login*
				FROM "order" o*
				LEFT JOIN users u ON u.id = o.user_id*
				LEFT JOIN status s ON s.id = o.status_id*
				WHERE number = \$1`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			dbr := Repository{Pool: poolMock}

			var rows pgx.Rows
			poolMock.EXPECT().
				Query(
					context.Background(),
					tt.query,
					tt.number).
				Return(rows, nil)

			_, err := dbr.QueryContext(context.Background(), tt.query, tt.number)
			assert.NoError(t, err)
		})
	}
}

func TestRepository_QueryRowContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name  string
		query string
		login string
	}{
		{
			name:  "test 1",
			query: `SELECT id, login, password, name FROM users WHERE login = $1`,
			login: "ramil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			poolMock := repository2.NewMockPooler(ctrl)
			dbr := Repository{Pool: poolMock}

			var row pgx.Row
			poolMock.EXPECT().
				QueryRow(
					context.Background(),
					tt.query,
					tt.login).
				Return(row)
			got := dbr.QueryRowContext(context.Background(), tt.query, tt.login)
			assert.Equal(t, got, row)
		})
	}
}
