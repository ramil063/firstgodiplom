package balance

// Balance описывает структуру баланса
type Balance struct {
	Current   float64 `json:"current"`   // Текущее
	Withdrawn int     `json:"withdrawn"` // Списания
}
