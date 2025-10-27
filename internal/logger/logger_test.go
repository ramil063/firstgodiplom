package logger

import "testing"

func TestWriteInfoLog(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteInfoLog("test 1")
		})
	}
}

func TestWriteDebugLog(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteDebugLog("test 1")
		})
	}
}

func TestWriteErrorLog(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteErrorLog("test 1")
		})
	}
}
