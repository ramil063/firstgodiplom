package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckGetMethodMw(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "test 1",
			method:       http.MethodPost,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "test 2",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "test 3",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			request := httptest.NewRequest(tc.method, "/", bytes.NewReader(body))
			// создаём новый Recorder
			w := httptest.NewRecorder()

			handlerToTest := CheckGetMethodMw(handler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tc.expectedCode, res.StatusCode, "Response code didn't match expected")

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			respBody, err := io.ReadAll(res.Body)
			require.NoError(t, err, "error making HTTP request")

			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(respBody))
			}
		})
	}
}

func TestCheckPostMethodMw(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	testCases := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "test 1",
			method:       http.MethodGet,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: "",
		},
		{
			name:         "test 2",
			method:       http.MethodPost,
			body:         "11114",
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "test 3",
			method:       http.MethodPost,
			body:         "11114",
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			request := httptest.NewRequest(tc.method, "/", bytes.NewReader(body))
			// создаём новый Recorder
			w := httptest.NewRecorder()

			handlerToTest := CheckPostMethodMw(handler)
			handlerToTest.ServeHTTP(w, request)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, tc.expectedCode, res.StatusCode, "Response code didn't match expected")

			// получаем и проверяем тело запроса
			defer res.Body.Close()
			respBody, err := io.ReadAll(res.Body)
			require.NoError(t, err, "error making HTTP request")

			// проверяем корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, string(respBody))
			}
		})
	}
}
