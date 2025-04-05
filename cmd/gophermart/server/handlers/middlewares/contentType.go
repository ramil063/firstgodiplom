package middlewares

import (
	"net/http"

	"github.com/ramil063/firstgodiplom/internal/logger"
)

// CheckContentTypeMiddleware проверка типа контента авторизации
func CheckContentTypeMiddleware(needContentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			contentType := r.Header.Get("Content-Type")

			if contentType != needContentType {
				logger.WriteErrorLog("wrong content type: " + contentType)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
