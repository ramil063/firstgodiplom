package balance

// Withdraw описывает структуру списания баллов
type Withdraw struct {
	OrderNumber string `json:"order"`                  // Заказ
	Sum         int    `json:"sum"`                    // Сумма
	ProcessedAt string `json:"processed_at,omitempty"` // Время
}
