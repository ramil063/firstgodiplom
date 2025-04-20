package errors

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDBError_Error(t *testing.T) {
	testTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	testErr := errors.New("test error")

	dbErr := &DBError{
		Time: testTime,
		Err:  testErr,
	}

	got := dbErr.Error()
	expected := "2024-01-01 12:00:00 test error"
	assert.Equal(t, expected, got)
}

func TestDBError_Unwrap(t *testing.T) {
	testErr := errors.New("test error")
	dbErr := &DBError{
		Time: time.Now(),
		Err:  testErr,
	}
	got := dbErr.Unwrap()
	assert.Equal(t, testErr, got)
}

func TestNewDBError(t *testing.T) {
	testErr := errors.New("test error")
	dbErr := NewDBError(testErr)

	if _, ok := dbErr.(*DBError); !ok {
		t.Error("NewDBError() did not return *DBError")
	}

	if !errors.Is(dbErr, testErr) {
		t.Error("NewDBError() wrapped error not accessible via errors.Is")
	}
}
