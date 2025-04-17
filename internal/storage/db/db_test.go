package db

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	db "github.com/ramil063/firstgodiplom/internal/storage/db/mocks"
)

func TestCheckPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
	}{
		{
			name: "test 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbr := db.NewMockDataBaser(ctrl)
			dbr.EXPECT().PingContext(gomock.Any()).Return(nil)
			err := CheckPing(dbr)
			assert.NoError(t, err)
		})
	}
}

func TestCreateTables(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
	}{
		{
			name: "test 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbr := db.NewMockDataBaser(ctrl)
			dbr.EXPECT().
				ExecContext(context.Background(), gomock.Any()).
				Return(nil, nil)

			err := CreateTables(dbr)
			assert.NoError(t, err)
		})
	}
}

func TestInit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name string
	}{
		{
			name: "test 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbr := db.NewMockDataBaser(ctrl)
			dbr.EXPECT().
				PingContext(gomock.Any()).
				Return(nil)
			dbr.EXPECT().
				ExecContext(context.Background(), gomock.Any()).
				Return(nil, nil)
			err := Init(dbr)
			assert.NoError(t, err)
		})
	}
}
