package errors

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewRequestError(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		statusCode int
		wantErr    string
	}{
		{
			name:       "test error with status 400",
			status:     "bad request",
			statusCode: 400,
			wantErr:    "bad request",
		},
		{
			name:       "test error with status 500",
			status:     "internal server error",
			statusCode: 500,
			wantErr:    "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewRequestError(tt.status, tt.statusCode)
			reqErr, ok := err.(*RequestError)
			assert.True(t, ok)
			assert.Equal(t, tt.statusCode, reqErr.StatusCode)
			assert.Equal(t, tt.wantErr, reqErr.Err.Error())
		})
	}
}

func TestRequestError_Error(t *testing.T) {
	tests := []struct {
		name       string
		time       time.Time
		statusCode int
		err        error
		want       string
	}{
		{
			name:       "test error string format",
			time:       time.Date(2024, 3, 15, 12, 0, 0, 0, time.UTC),
			statusCode: 400,
			err:        errors.New("test error"),
			want:       "2024-03-15 12:00:00 Status:400 Error:test error",
		},
		{
			name:       "test another error",
			time:       time.Date(2024, 3, 15, 15, 30, 0, 0, time.UTC),
			statusCode: 500,
			err:        errors.New("another error"),
			want:       "2024-03-15 15:30:00 Status:500 Error:another error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &RequestError{
				Time:       tt.time,
				StatusCode: tt.statusCode,
				Err:        tt.err,
			}
			assert.Equal(t, tt.want, e.Error())
		})
	}
}

func TestRequestError_Unwrap(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantErr error
	}{
		{
			name:    "test unwrap error",
			err:     errors.New("original error"),
			wantErr: errors.New("original error"),
		},
		{
			name:    "test unwrap nil error",
			err:     nil,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &RequestError{
				Time:       time.Now(),
				StatusCode: 400,
				Err:        tt.err,
			}
			if tt.wantErr == nil {
				assert.Nil(t, e.Unwrap())
			} else {
				assert.Equal(t, tt.wantErr.Error(), e.Unwrap().Error())
			}
		})
	}
}
