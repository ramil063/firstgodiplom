package balance

// Balance описывает структуру баланса
type Balance struct {
	ID        int     `json:"id,omitempty"` // Идентификатор
	Current   float32 `json:"current"`      // Текущее
	Withdrawn float32 `json:"withdrawn"`    // Списания
}
