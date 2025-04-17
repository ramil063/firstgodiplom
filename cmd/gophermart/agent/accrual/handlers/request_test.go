package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ClientMock struct {
	httpClient MockHTTPClient
}
type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func (c ClientMock) NewRequest(method string, url string, body []byte) (*http.Request, error) {
	return httptest.NewRequest(method, url, bytes.NewReader(body)), nil
}

func (c ClientMock) SendRequest(method string, url string, body []byte) (int, io.ReadCloser, http.Header, error) {
	req, err := c.NewRequest(method, url, body)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	return res.StatusCode, res.Body, res.Header, nil
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"test 1", "handlers.client"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, reflect.ValueOf(NewClient()).Type().String(), "NewClient()")
		})
	}
}

func Test_client_NewRequest(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"test 1", "*http.Request"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient()
			req, err := c.NewRequest("GET", "/", []byte{})
			assert.Equalf(t, tt.want, reflect.ValueOf(req).Type().String(), "NewRequest()")
			assert.NoError(t, err, "NewRequest() no error")
		})
	}
}

func Test_client_SendRequest(t *testing.T) {
	mockHttpClient := &MockHTTPClient{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       nil,
				Header:     http.Header{},
			}, nil
		},
	}

	type args struct {
		method string
		url    string
		body   []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		want1   io.ReadCloser
		want2   http.Header
		wantErr bool
	}{
		{
			"test1",
			args{"GET", "http://localhost:8080", []byte("")},
			http.StatusOK,
			nil,
			http.Header{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := ClientMock{
				httpClient: *mockHttpClient,
			}
			got, got1, got2, err := c.SendRequest(tt.args.method, tt.args.url, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
			assert.Equal(t, got2, tt.want2)
			assert.NoError(t, err)
		})
	}
}
