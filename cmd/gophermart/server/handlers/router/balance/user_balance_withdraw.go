package balance

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/theplant/luhn"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	balanceData "github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
	internalContextKeys "github.com/ramil063/firstgodiplom/internal/constants/context"
	internalErrors "github.com/ramil063/firstgodiplom/internal/errors"
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

	tokenData, ok := r.Context().Value(internalContextKeys.AccessTokenData).(user.AccessTokenData)
	if !ok {
		logger.WriteErrorLog("putOrder AccessTokenData not found in context")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = dbs.AddWithdrawFromBalance(withdraw, tokenData.Login)
	if errors.Is(err, internalErrors.ErrNotEnoughBalance) {
		logger.WriteErrorLog("not enough balance")
		rw.WriteHeader(http.StatusPaymentRequired)
		return
	}
	if err != nil {
		logger.WriteErrorLog("AddWithdraw AddWithdraw error:" + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.WriteInfoLog(tokenData.Login + "withdraw added successfully")
	rw.WriteHeader(http.StatusOK)
}
