package balance

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
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

func TestAddWithdraw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		wantCode    int
		login       string
		accessToken string
		orderNumber string
		sumWithdraw float32
		balance     balance.Balance
	}{
		{
			name:        "test 1",
			wantCode:    http.StatusOK,
			login:       "ramil",
			accessToken: "123",
			orderNumber: "603414218808776",
			sumWithdraw: 655.66,
			balance: balance.Balance{
				Current: 700,
			},
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
			withdrawMock := balance.Withdraw{
				OrderNumber: tt.orderNumber,
				Sum:         tt.sumWithdraw,
			}
			storageMock.EXPECT().AddWithdrawFromBalance(withdrawMock, tt.login).Return(nil)

			sumStr := strconv.FormatFloat(float64(tt.sumWithdraw), 'f', -1, 32)
			body := []byte("{\n    \"order\": \"" + tt.orderNumber + "\",\n    \"sum\": " + sumStr + "\n}")
			request := httptest.NewRequest("POST", "/api/user/balance/withdraw", bytes.NewReader(body))
			request.Header.Set("Authorization", "Bearer "+tt.accessToken)

			// создаём новый Recorder
			w := httptest.NewRecorder()

			addWithdrawHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
				AddWithdraw(rw, r, storageMock)
			}
			handlerToTest := middlewares.CheckAuthMiddleware(storageMock)(http.HandlerFunc(addWithdrawHandlerFunction))
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.wantCode, res.StatusCode)

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}
