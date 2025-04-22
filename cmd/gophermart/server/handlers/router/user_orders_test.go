package router

import (
	"bytes"
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
)

func Test_getOrders(t *testing.T) {
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
			ordersMock := []user.Order{{}}
			storageMock.EXPECT().GetOrders(tt.login).Return(ordersMock, nil)

			request := httptest.NewRequest("GET", "/api/user/orders", nil)
			request.Header.Set("Authorization", "Bearer "+tt.accessToken)

			// создаём новый Recorder
			w := httptest.NewRecorder()

			getOrderHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
				getOrders(rw, r, storageMock)
			}
			handlerToTest := middlewares.CheckAuthMiddleware(storageMock)(http.HandlerFunc(getOrderHandlerFunction))
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.wantCode, res.StatusCode)

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}

func Test_putOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name        string
		wantCode    int
		login       string
		accessToken string
		orderNumber string
	}{
		{
			name:        "test 1",
			wantCode:    http.StatusAccepted,
			login:       "ramil",
			accessToken: "123",
			orderNumber: "2850230465763",
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
			ordersMock := user.Order{}
			storageMock.EXPECT().GetOrder(tt.orderNumber).Return(ordersMock, nil)
			accessTokenData := user.AccessTokenData{
				Login:                tt.login,
				AccessToken:          tt.accessToken,
				AccessTokenExpiredAt: time.Now().Add(time.Minute).Unix(),
			}
			storageMock.EXPECT().AddOrder(tt.orderNumber, accessTokenData).Return(nil)

			body := []byte(tt.orderNumber)
			request := httptest.NewRequest("POST", "/api/user/orders", bytes.NewReader(body))
			request.Header.Set("Authorization", "Bearer "+tt.accessToken)

			// создаём новый Recorder
			w := httptest.NewRecorder()

			putOrderHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
				putOrder(rw, r, storageMock)
			}
			handlerToTest := middlewares.CheckAuthMiddleware(storageMock)(http.HandlerFunc(putOrderHandlerFunction))
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.wantCode, res.StatusCode)

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}
