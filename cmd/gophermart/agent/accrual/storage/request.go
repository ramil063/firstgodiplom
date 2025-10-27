package storage

// DefaultRetryAfterInterval Интервал в секундах через который нужно опросить сервис акруал
// если он вернул специальный статус, но не вернул через сколько его опросить
var DefaultRetryAfterInterval = 60
