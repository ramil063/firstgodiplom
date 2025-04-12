package router

import (
	"encoding/json"
	"net/http"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

// getWithdrawals получение списаний баллов с баланса
func getWithdrawals(rw http.ResponseWriter, r *http.Request, dbs storage.Storager) {
	token := auth.GetTokenFromHeader(r)
	tokenData, err := dbs.GetAccessTokenData(token)
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	withdrawals, err := dbs.GetWithdrawals(tokenData.Login)
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(withdrawals) == 0 {
		logger.WriteErrorLog("no rows returned")
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	logger.WriteInfoLog("orders show")

	enc := json.NewEncoder(rw)
	if err = enc.Encode(withdrawals); err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
