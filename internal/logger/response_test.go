package logger

import (
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseLogger(t *testing.T) {
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

			handlerToTest := ResponseLogger(nextHandler)
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

func Test_loggingResponseWriter_Write(t *testing.T) {
	type fields struct {
		ResponseWriter http.ResponseWriter
		responseData   *responseData
	}
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		newSize int
		wantErr bool
	}{
		{
			"test 1",
			fields{
				ResponseWriter: httptest.NewRecorder(),
				responseData:   &responseData{status: 200, size: 10},
			},
			args{b: big.NewInt(10).Bytes()},
			1,
			11,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &loggingResponseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				responseData:   tt.fields.responseData,
			}
			got, err := r.Write(tt.args.b)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.newSize, r.responseData.size)
		})
	}
}

func Test_loggingResponseWriter_WriteHeader(t *testing.T) {
	type fields struct {
		ResponseWriter http.ResponseWriter
		responseData   *responseData
	}
	type args struct {
		statusCode int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			"test 1",
			fields{
				ResponseWriter: httptest.NewRecorder(),
				responseData:   &responseData{status: 200, size: 10},
			},
			args{statusCode: 200},
			1,
			false,
		},
		{
			"test 1",
			fields{
				ResponseWriter: httptest.NewRecorder(),
				responseData:   &responseData{status: 404, size: 10},
			},
			args{statusCode: 404},
			1,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &loggingResponseWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				responseData:   tt.fields.responseData,
			}
			r.WriteHeader(tt.args.statusCode)
			assert.Equal(t, tt.args.statusCode, r.responseData.status)
		})
	}
}
