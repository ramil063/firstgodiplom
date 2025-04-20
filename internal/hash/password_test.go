package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPasswordHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "correct password",
			password: "test123",
			want:     true,
		},
		{
			name:     "incorrect password",
			password: "wrongpass",
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate hash for the test case
			if tt.hash == "" {
				var err error
				tt.hash, err = GetPasswordHash(tt.password)
				assert.NoError(t, err)
			}

			if tt.want {
				assert.True(t, CheckPasswordHash(tt.password, tt.hash))
			} else {
				assert.False(t, CheckPasswordHash("wrongpassword", tt.hash))
			}
		})
	}
}

func TestGetPasswordHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "test123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false,
		},
		{
			name:     "long password",
			password: "verylongpasswordtestingmorethan72bytes_verylongpasswordtestingmorethan72bytes",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GetPasswordHash(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				// Verify the hash works with CheckPasswordHash
				assert.True(t, CheckPasswordHash(tt.password, hash))
			}
		})
	}
}
