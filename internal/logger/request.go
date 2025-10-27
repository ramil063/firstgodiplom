package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize() error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	// устанавливаем синглтон
	Log = zl
	return nil
}

// RequestLogger — middleware-логер для входящих HTTP-запросов.
func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// функция Now() возвращает текущее время
		start := time.Now()
		h.ServeHTTP(w, r)
		// Since возвращает разницу во времени между start
		// и моментом вызова Since. Таким образом можно посчитать
		// время выполнения запроса.
		duration := time.Since(start)
		Log.Info("got incoming HTTP request",
			zap.String("URI", r.URL.RequestURI()),
			zap.String("method", r.Method),
			zap.String("duration", duration.String()),
		)
	})
}
