package router

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ramil063/firstgodiplom/internal/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	storage2 "github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/mocks"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
)

func Test_userLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		wantCode int
		login    string
		password string
	}{
		{
			name:     "test 1",
			wantCode: http.StatusOK,
			login:    "ramil",
			password: "123456",
		},
	}
	for _, tt := range tests {

		storageMock := storage2.NewMockStorager(ctrl)
		passwordHash, err := hash.GetPasswordHash(tt.password)
		assert.NoError(t, err)

		userMock := user.User{
			Login:        tt.login,
			PasswordHash: passwordHash,
		}
		storageMock.EXPECT().GetUser(context.Background(), tt.login).Return(userMock, nil)
		storageMock.EXPECT().UpdateToken(context.Background(), "ramil", gomock.Any(), gomock.Any()).Return(nil)

		body := []byte("{\n    \"login\": \"" + tt.login + "\",\n    \"password\": \"" + tt.password + "\"\n}")
		request := httptest.NewRequest("POST", "/api/user/login", bytes.NewReader(body))

		// создаём новый Recorder
		w := httptest.NewRecorder()

		userLoginHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
			userLogin(rw, r, storageMock)
		}
		handlerToTest := http.HandlerFunc(userLoginHandlerFunction)
		handlerToTest.ServeHTTP(w, request)

		res := w.Result()

		assert.Equal(t, tt.wantCode, res.StatusCode)

		defer res.Body.Close()
		_, err = io.ReadAll(res.Body)
		require.NoError(t, err)
	}
}
