package middlewares

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

	storage2 "github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/mocks"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
)

func TestCheckAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		want     func(http.Handler) http.Handler
		wantCode int
	}{
		{
			name: "test 1",
			want: func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/1", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			request.Header.Set("Authorization", "Bearer 123")

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

			storageMock := storage2.NewMockStorager(ctrl)
			token := user.AccessTokenData{
				AccessToken:          "123",
				AccessTokenExpiredAt: time.Now().Add(time.Minute).Unix(),
			}
			storageMock.EXPECT().GetAccessTokenData(context.Background(), "123").Return(token, nil)

			handlerToTest := CheckAuthMiddleware(storageMock)(handler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.wantCode, res.StatusCode)

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}
