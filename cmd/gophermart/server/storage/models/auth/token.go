package auth

var TokenExpiredSeconds int64 = 60 * 60 * 24 * 30

// Token описывает токен авторизации
type Token struct {
	Token string `json:"token"` // Токен
}
