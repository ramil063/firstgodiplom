package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	storage2 "github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/mocks"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
)

func Test_getWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		wantCode    int
		login       string
		accessToken string
	}{
		{
			name:        "test 1",
			wantCode:    http.StatusOK,
			login:       "ramil",
			accessToken: "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storageMock := storage2.NewMockStorager(ctrl)
			token := user.AccessTokenData{
				Login:                tt.login,
				AccessToken:          tt.accessToken,
				AccessTokenExpiredAt: time.Now().Add(time.Minute).Unix(),
			}
			storageMock.EXPECT().GetAccessTokenData(tt.accessToken).Return(token, nil)
			withdrawalsMock := []balance.Withdraw{{}}
			storageMock.EXPECT().GetWithdrawals(tt.login).Return(withdrawalsMock, nil)

			request := httptest.NewRequest("GET", "/api/user/withdrawals", nil)
			request.Header.Set("Authorization", "Bearer "+tt.accessToken)

			// создаём новый Recorder
			w := httptest.NewRecorder()

			getWithdrawHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
				getWithdrawals(rw, r, storageMock)
			}
			handlerToTest := http.HandlerFunc(getWithdrawHandlerFunction)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.wantCode, res.StatusCode)

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}
