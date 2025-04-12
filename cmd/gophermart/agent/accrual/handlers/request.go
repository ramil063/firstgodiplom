package handlers

import (
	"bytes"
	"io"
	"net/http"
)

type Clienter interface {
	NewRequest(method string, url string) (*http.Request, error)
	SendRequest(method string, url string, body []byte) (int, io.ReadCloser, http.Header, error)
}

type client struct {
	httpClient *http.Client
}

// NewClient создание нового клиента
func NewClient() Clienter {
	return client{
		httpClient: &http.Client{},
	}
}

// SendRequest отправка запроса
func (c client) SendRequest(method string, url string, body []byte) (int, io.ReadCloser, http.Header, error) {

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, nil, err
	}

	return res.StatusCode, res.Body, res.Header, nil
}

// NewRequest создание нового запроса
func (c client) NewRequest(method string, url string) (*http.Request, error) {
	return http.NewRequest(method, url, nil)
}
