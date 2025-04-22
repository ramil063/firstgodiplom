package balance

import (
	"encoding/json"
	"net/http"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	internalContextKeys "github.com/ramil063/firstgodiplom/internal/constants/context"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

// GetBalance получение баланса пользователя
func GetBalance(rw http.ResponseWriter, r *http.Request, s storage.Storager) {

	tokenData, ok := r.Context().Value(internalContextKeys.AccessTokenData).(user.AccessTokenData)
	if !ok {
		logger.WriteErrorLog("putOrder AccessTokenData not found in context")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance, err := s.GetBalance(tokenData.Login)
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	logger.WriteInfoLog("balance show")

	enc := json.NewEncoder(rw)
	if err = enc.Encode(balance); err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
