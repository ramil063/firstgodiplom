package user

// Order описывает структуру заказа
type Order struct {
	Id         int    `json:"id,omitempty"`   // Идентификатор
	Number     string `json:"number"`         // Номер
	Status     string `json:"status"`         // Статус
	Accrual    int    `json:"accrual"`        // Начисления
	UploadedAt string `json:"uploaded_at"`    // Дата загрузки
	UserLogin  string `json:"user,omitempty"` // Пользователь
}

// OrderCheckAccrual описывает структуру заказа
type OrderCheckAccrual struct {
	Id        int    `json:"id,omitempty"`   // Идентификатор
	Number    string `json:"number"`         // Номер
	Status    string `json:"status"`         // Статус
	Accrual   int    `json:"accrual"`        // Начисления
	UserLogin string `json:"user,omitempty"` // Пользователь
}
