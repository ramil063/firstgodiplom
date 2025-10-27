package flags

import "testing"

func TestParseFlags(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "test 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ParseFlags()
		})
	}
}
