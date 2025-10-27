package router

import (
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	storage2 "github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/mocks"
)

func TestRouter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
		want chi.Router
	}{
		{
			name: "test 1",
			want: chi.NewRouter(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageMock := storage2.NewMockStorager(ctrl)
			got := Router(storageMock)
			assert.Equal(t, reflect.ValueOf(tt.want).Type().String(), reflect.ValueOf(got).Type().String())
		})
	}
}
