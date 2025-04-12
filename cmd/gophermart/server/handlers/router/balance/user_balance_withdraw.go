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

// AddWithdraw добавление списания средств с баланса
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
	logger.WriteInfoLog("AddWithdraw:" + string(logMsg))

	if withdraw.Sum <= 0 {
		logger.WriteErrorLog("AddWithdraw sum must be positive")
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
		logger.WriteErrorLog("AddWithdraw wrong format luhn number:" + withdraw.OrderNumber)
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	token := auth.GetTokenFromHeader(r)
	tokenData, err := dbs.GetAccessTokenData(token)
	if err != nil {
		logger.WriteErrorLog("AddWithdraw GetAccessTokenData error:" + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	balance, err := dbs.GetBalance(tokenData.Login)
	if err != nil {
		logger.WriteErrorLog("AddWithdraw GetBalance error:" + err.Error())
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
		logger.WriteErrorLog("AddWithdraw AddWithdraw error:" + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.WriteInfoLog(tokenData.Login + "withdraw added successfully")
	rw.WriteHeader(http.StatusOK)
}
