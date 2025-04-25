package balance

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/middlewares"
	storage2 "github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/mocks"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
)

func TestGetBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		wantCode int
	}{
		{
			name:     "test 1",
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storageMock := storage2.NewMockStorager(ctrl)
			token := user.AccessTokenData{
				Login:                "ramil",
				AccessToken:          "123",
				AccessTokenExpiredAt: time.Now().Add(time.Minute).Unix(),
			}
			storageMock.EXPECT().GetAccessTokenData(context.Background(), "123").Return(token, nil)
			balanceMock := balance.Balance{}
			storageMock.EXPECT().GetBalance(gomock.Any(), "ramil").Return(balanceMock, nil)

			request := httptest.NewRequest("GET", "/api/user/balance", nil)
			request.Header.Set("Authorization", "Bearer 123")

			// создаём новый Recorder
			w := httptest.NewRecorder()

			getBalanceHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
				GetBalance(rw, r, storageMock)
			}
			handlerToTest := middlewares.CheckAuthMiddleware(storageMock)(http.HandlerFunc(getBalanceHandlerFunction))
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.wantCode, res.StatusCode)

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}
