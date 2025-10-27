package user

// Login описывает входную структуру при логине пользователя
type Login struct {
	Login    string `json:"login"`    // Логин
	Password string `json:"password"` // Пароль
}
