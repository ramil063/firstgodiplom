package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	internalErrors "github.com/ramil063/firstgodiplom/internal/errors"
	"github.com/ramil063/firstgodiplom/internal/hash"
)

// AuthenticateUser аутентифицировать пользователя
func AuthenticateUser(s storage.Storager, login string) (auth.Token, error) {
	var err error
	var t auth.Token

	t.Token, err = hash.RandomHex(20)
	if err != nil {
		return t, err
	}

	expiredAt := time.Now().Unix() + auth.TokenExpiredSeconds
	err = s.UpdateToken(login, t, expiredAt)
	return t, err
}

// CheckAuthenticatedUser проверка пользователя на аутентификацию
func CheckAuthenticatedUser(s storage.Storager, token string) (user.AccessTokenData, error) {
	var err error
	var tokenData user.AccessTokenData
	tokenData, err = s.GetAccessTokenData(token)

	if errors.Is(err, pgx.ErrNoRows) {
		return tokenData, internalErrors.ErrIncorrectToken
	}

	if tokenData.AccessTokenExpiredAt < time.Now().Unix() {
		return tokenData, internalErrors.ErrExpiredToken
	}
	return tokenData, err
}

// GetTokenFromHeader получить токен из хедера
func GetTokenFromHeader(r *http.Request) string {
	authorization := r.Header.Get("Authorization")
	return strings.TrimSpace(strings.Replace(authorization, "Bearer", "", 1))
}
