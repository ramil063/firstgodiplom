package user

// Register описывает входную структуру при регистрации пользователя
type Register struct {
	Login    string `json:"login"`          // Логин
	Password string `json:"password"`       // Пароль
	Name     string `json:"name,omitempty"` // Имя
}
