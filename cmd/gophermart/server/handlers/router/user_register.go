package router

import (
	"encoding/json"
	"errors"
	"net/http"

	authHandlers "github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	internalErrors "github.com/ramil063/firstgodiplom/internal/errors"
	"github.com/ramil063/firstgodiplom/internal/hash"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

// userRegister регистрация пользователя
func userRegister(rw http.ResponseWriter, r *http.Request, s storage.Storager) {
	var register user.Register
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&register)

	if err != nil {
		logger.WriteDebugLog(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	logMsg, _ := json.Marshal(register)
	logger.WriteInfoLog(string(logMsg))

	passwordHash, err := hash.GetPasswordHash(register.Password)
	if err != nil {
		logger.WriteErrorLog("error in crypt password")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(register.Login) < 1 || len(register.Password) < 1 {
		logger.WriteErrorLog("empty login/password")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.AddUserData(register, passwordHash)

	if errors.Is(err, internalErrors.ErrUniqueViolation) {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusConflict)
		return
	}
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var t auth.Token
	t, err = authHandlers.AuthenticateUser(s, register.Login)
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

	logger.WriteDebugLog("sending HTTP 200 response")
}
