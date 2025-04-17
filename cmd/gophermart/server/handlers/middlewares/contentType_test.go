package middlewares

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckContentTypeMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		needContentType string
		wantCode        int
	}{
		{
			name:            "test 1",
			needContentType: "application/json",
			wantCode:        http.StatusOK,
		},
		{
			name:            "test 2",
			needContentType: "text/plain",
			wantCode:        http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/1", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			request.Header.Set("Content-Type", "application/json")

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

			handlerToTest := CheckContentTypeMiddleware(tt.needContentType)(handler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()

			assert.Equal(t, tt.wantCode, res.StatusCode)

			defer res.Body.Close()
			_, err := io.ReadAll(res.Body)
			require.NoError(t, err)
		})
	}
}
