package balance

import (
	"encoding/json"
	"net/http"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

func GetBalance(rw http.ResponseWriter, r *http.Request, dbs storage.Storager) {

	token := auth.GetTokenFromHeader(r)
	tokenData, err := dbs.GetAccessTokenData(token)
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance, err := dbs.GetBalance(tokenData.Login)
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
