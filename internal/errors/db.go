package errors

import (
	"errors"
	"fmt"
	"time"
)

var ErrUniqueViolation = errors.New("unique violation")

type DBError struct {
	Time time.Time
	Err  error
}

func (e *DBError) Error() string {
	return fmt.Sprintf("%v %v", e.Time.Format("2006-01-02 15:04:05"), e.Err)
}

func (e *DBError) Unwrap() error {
	return e.Err
}

// NewDBError записывает ошибку err в тип DBError c текущим временем.
func NewDBError(err error) error {
	return &DBError{
		Time: time.Now(),
		Err:  err,
	}
}
