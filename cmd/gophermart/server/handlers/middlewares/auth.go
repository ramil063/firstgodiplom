package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	intenalContextKeys "github.com/ramil063/firstgodiplom/internal/constants/context"
	internalErrors "github.com/ramil063/firstgodiplom/internal/errors"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

// CheckAuthMiddleware проверка токена авторизации
func CheckAuthMiddleware(s storage.Storager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := auth.GetTokenFromHeader(r)

			tokenData, err := auth.CheckAuthenticatedUser(r.Context(), s, token)

			if errors.Is(err, internalErrors.ErrIncorrectToken) || errors.Is(err, internalErrors.ErrExpiredToken) {
				logger.WriteErrorLog("token is incorrect:" + token)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if err != nil {
				logger.WriteErrorLog(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			k := intenalContextKeys.AccessTokenData
			ctx := context.WithValue(r.Context(), k, tokenData)
			newRequest := r.WithContext(ctx)

			next.ServeHTTP(w, newRequest)
		})
	}
}
