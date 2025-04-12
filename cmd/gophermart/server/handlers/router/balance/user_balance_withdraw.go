package balance

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/theplant/luhn"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	balanceData "github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

func AddWithdraw(rw http.ResponseWriter, r *http.Request, dbs storage.Storager) {
	var withdraw balanceData.Withdraw
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&withdraw)

	if err != nil {
		logger.WriteDebugLog(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	logMsg, _ := json.Marshal(withdraw)
	logger.WriteInfoLog(string(logMsg))

	if withdraw.Sum <= 0 {
		logger.WriteErrorLog("sum must be positive")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	num, err := strconv.Atoi(withdraw.OrderNumber)
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !luhn.Valid(num) {
		logger.WriteErrorLog("wrong format luhn number")
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

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

	if balance.Current < withdraw.Sum {
		logger.WriteErrorLog("not enough balance")
		rw.WriteHeader(http.StatusPaymentRequired)
		return
	}

	err = dbs.AddWithdraw(withdraw, tokenData.Login)
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.WriteInfoLog("withdraw added successfully")
	rw.WriteHeader(http.StatusOK)
}
