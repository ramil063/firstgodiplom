package user

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository"
	repository2 "github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository/mocks"
)

func TestAddUser(t *testing.T) {

	tests := []struct {
		name     string
		login    string
		password string
		userName string
	}{
		{
			name:     "test 1",
			login:    "ramil",
			password: "password",
			userName: "ramil",
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
			mock.ExpectExec(`^INSERT INTO users \(login, password, name\) VALUES \(\$1, \$2, \$3\)$`).
				WithArgs(
					tt.login,
					tt.password,
					tt.userName).
				WillReturnResult(expectedCommandTag)
			mock.ExpectCommit()

			tx, err := mock.Begin(context.Background())
			if err != nil {
				t.Fatal(err)
			}

			_, err = AddUser(tx, tt.login, tt.password, tt.userName)
			assert.NoError(t, err)
		})
	}
}

func TestUpdateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name      string
		login     string
		token     string
		expiredAt int64
	}{
		{
			name:      "test 1",
			login:     "ramil",
			token:     "token",
			expiredAt: 10,
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
					`UPDATE users SET access_token = $1, access_token_expired_at = $2 WHERE login = $3`,
					tt.token,
					tt.expiredAt,
					tt.login).
				Return(expectedCommandTag, nil)
			_, err := UpdateToken(&dbr, tt.login, tt.token, tt.expiredAt)
			assert.NoError(t, err)
		})
	}
}
