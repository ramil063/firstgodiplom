package middlewares

import (
	"net/http"
	"strings"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/gzip"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

// GZIPMiddleware нужен для сжатия входящих и выходных данных
func GZIPMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// по умолчанию устанавливаем оригинальный http.ResponseWriter как тот,
		// который будем передавать следующей функции
		ow := w

		// проверяем, что клиент умеет получать от сервера сжатые данные в формате gzip
		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		contentType := r.Header.Get("Content-Type")
		accept := r.Header.Get("Accept")
		applicationJSON := strings.Contains(contentType, "application/json")
		textHTML := strings.Contains(contentType, "text/html")
		acceptTextHTML := strings.Contains(accept, "text/html")

		if supportsGzip && (applicationJSON || textHTML || acceptTextHTML) {
			// оборачиваем оригинальный http.ResponseWriter новым с поддержкой сжатия
			cw := gzip.NewCompressWriter(w)
			// меняем оригинальный http.ResponseWriter на новый
			ow = cw
			// не забываем отправить клиенту все сжатые данные после завершения middleware
			defer cw.Close()
		}

		// проверяем, что клиент отправил серверу сжатые данные в формате gzip
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip && (applicationJSON || textHTML || acceptTextHTML) {
			// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
			cr, err := gzip.NewCompressReader(r.Body)
			if err != nil {
				logger.WriteErrorLog(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			// меняем тело запроса на новое
			r.Body = cr
			defer cr.Close()
		}

		// передаём управление хендлеру
		next.ServeHTTP(ow, r)
	})
}
