package user

// User описывает пользователя
type User struct {
	ID           int    `json:"id,omitempty"`
	Login        string `json:"login"`          // Логин
	PasswordHash string `json:"password"`       // Пароль
	Name         string `json:"name,omitempty"` // Имя
}

// AccessTokenData описывает данные токена пользователя
type AccessTokenData struct {
	Login                string `json:"login"`                   // Логин
	AccessToken          string `json:"access_token"`            // Токен авторизации
	AccessTokenExpiredAt int64  `json:"access_token_expired_at"` // Время истечения токена авторизации
}
