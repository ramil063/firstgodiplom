package hash

import (
	"crypto/rand"
	"encoding/hex"
)

// RandomHex рандомный набор символов
func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
