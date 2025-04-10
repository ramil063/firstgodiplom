package storage

// Order описывает структуру заказа
type Order struct {
	Id         int    `json:"id,omitempty,omitempty"` // Идентификатор
	Order      string `json:"order"`                  // Номер заказа
	Status     string `json:"status"`                 // Статус
	Accrual    int    `json:"accrual"`                // Начисления
	UploadedAt string `json:"uploaded_at,omitempty"`  // Дата загрузки
	UserLogin  string `json:"user,omitempty"`         // Пользователь
}
