package logger

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{"test 1", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Initialize()
			assert.NoError(t, err)
		})
	}
}

func TestRequestLogger(t *testing.T) {
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name       string
		pathValues map[string]string
		want       want
	}{
		{"test 1", map[string]string{"type": "gauge", "metric": "a"}, want{http.StatusOK, ""}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/value", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()

			request.SetPathValue("type", test.pathValues["type"])
			if metric, ok := test.pathValues["metric"]; ok {
				request.SetPathValue("metric", metric)
			}

			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})

			handlerToTest := RequestLogger(nextHandler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, test.want.code, res.StatusCode)
			// получаем и проверяем тело запроса

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
		})
	}
}
