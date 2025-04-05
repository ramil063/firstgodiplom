package middlewares

import (
	"net/http"

	"github.com/ramil063/firstgodiplom/internal/logger"
)

// CheckMethodMw middleware для проверки метода запроса
func CheckMethodMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// разрешаем только POST, GET запросы
		if r.Method != http.MethodPost && r.Method != http.MethodGet {
			logger.WriteDebugLog("Incorrect method")
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("Incorrect method"))
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CheckPostMethodMw middleware для проверки метода запроса
func CheckPostMethodMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// разрешаем только POST запросы
		if r.Method != http.MethodPost {
			logger.WriteDebugLog("got request with bad method:" + r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CheckGetMethodMw middleware для проверки метода запроса
func CheckGetMethodMw(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// разрешаем только GET запросы
		if r.Method != http.MethodGet {
			logger.WriteDebugLog("got request with bad method:" + r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}
