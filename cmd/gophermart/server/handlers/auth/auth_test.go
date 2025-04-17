package auth

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	storage2 "github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/mocks"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
)

func TestAuthenticateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		want string
	}{
		{
			name: "test 1",
			want: "auth.Token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			storageMock := storage2.NewMockStorager(ctrl)
			storageMock.EXPECT().UpdateToken("ramil", gomock.Any(), gomock.Any()).Return(nil)

			got, err := AuthenticateUser(storageMock, "ramil")
			assert.Equal(t, tt.want, reflect.ValueOf(got).Type().String())
			assert.NoError(t, err)
		})
	}
}

func TestCheckAuthenticatedUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name  string
		token string
	}{
		{"test 1", "123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var td user.AccessTokenData
			td.AccessTokenExpiredAt = time.Now().Add(time.Minute).Unix()
			storageMock := storage2.NewMockStorager(ctrl)
			storageMock.EXPECT().GetAccessTokenData(tt.token).Return(td, nil)

			err := CheckAuthenticatedUser(storageMock, tt.token)
			assert.NoError(t, err)
		})
	}
}

func TestGetTokenFromHeader(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	tests := []struct {
		name    string
		request *http.Request
		want    string
	}{
		{
			name:    "test 1",
			request: req,
			want:    "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.request.Header.Set("Authorization", "Bearer 123")
			got := GetTokenFromHeader(tt.request)
			assert.Equal(t, tt.want, got)
		})
	}
}
