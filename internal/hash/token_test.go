package hash

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomHex(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		wantLen int
		wantErr bool
	}{
		{
			name:    "test positive length",
			length:  8,
			wantLen: 16, // hex encoding doubles the length
			wantErr: false,
		},
		{
			name:    "test zero length",
			length:  0,
			wantLen: 0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandomHex(tt.length)
			assert.NoError(t, err)
			assert.Len(t, got, tt.wantLen)
			// Verify the output is valid hex
			_, err = hex.DecodeString(got)
			assert.NoError(t, err)
		})
	}
}
