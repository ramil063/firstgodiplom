package storage

// Order описывает структуру заказа
type Order struct {
	ID         int     `json:"id,omitempty"`          // Идентификатор
	Order      string  `json:"order"`                 // Номер заказа
	Status     string  `json:"status"`                // Статус
	Accrual    float32 `json:"accrual"`               // Начисления
	UploadedAt string  `json:"uploaded_at,omitempty"` // Дата загрузки
	UserLogin  string  `json:"user,omitempty"`        // Пользователь
}
