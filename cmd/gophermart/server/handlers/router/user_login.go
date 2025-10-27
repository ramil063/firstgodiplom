package router

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v4"

	authHandlers "github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/internal/hash"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

// userLogin авторизация пользователя
func userLogin(rw http.ResponseWriter, r *http.Request, dbs storage.Storager) {
	var login user.Login
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&login)

	if err != nil {
		logger.WriteDebugLog(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	logMsg, _ := json.Marshal(login)
	logger.WriteInfoLog(string(logMsg))

	u, err := dbs.GetUser(r.Context(), login.Login)

	if errors.Is(err, pgx.ErrNoRows) || !hash.CheckPasswordHash(login.Password, u.PasswordHash) {
		logger.WriteErrorLog("login/password incorrect")
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var t auth.Token
	t, err = authHandlers.AuthenticateUser(r.Context(), dbs, login.Login)
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Authorization", "Bearer "+t.Token)
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(rw)
	if err = enc.Encode(t); err != nil {
		logger.WriteErrorLog(err.Error())
		return
	}

	logger.WriteDebugLog("user login successfully:" + login.Login)
}
