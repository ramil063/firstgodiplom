package env

import (
	"os"
	"testing"
)

func TestInitEnvironmentVariables(t *testing.T) {
	// Save original environment variable
	originalAppEnv := os.Getenv("APP_ENV")
	defer os.Setenv("APP_ENV", originalAppEnv)

	tests := []struct {
		name     string
		envValue string
		want     string
	}{
		{
			name:     "empty environment",
			envValue: "",
			want:     "",
		},
		{
			name:     "development environment",
			envValue: "development",
			want:     "development",
		},
		{
			name:     "production environment",
			envValue: "production",
			want:     "production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set test environment variable
			os.Setenv("APP_ENV", tt.envValue)

			// Reset package variable
			AppEnv = ""

			InitEnvironmentVariables()

			if AppEnv != tt.want {
				t.Errorf("InitEnvironmentVariables() = %v, want %v", AppEnv, tt.want)
			}
		})
	}
}
