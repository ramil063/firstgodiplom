package errors

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

var IncorrectToken = errors.New("incorrect access token")
var ExpiredToken = errors.New("expired access token")

type RequestError struct {
	Time       time.Time
	StatusCode int
	Err        error
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("%v Status:%v Error:%v", e.Time.Format("2006-01-02 15:04:05"), strconv.Itoa(e.StatusCode), e.Err)
}

func (e *RequestError) Unwrap() error {
	return e.Err
}

// NewRequestError записывает ошибку err в тип RequestError c текущим временем и статусом.
func NewRequestError(status string, statusCode int) error {
	return &RequestError{
		Time:       time.Now(),
		StatusCode: statusCode,
		Err:        errors.New(status),
	}
}
