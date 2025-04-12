package router

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/theplant/luhn"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

// putOrder добавление нового заказа
func putOrder(rw http.ResponseWriter, r *http.Request, dbs storage.Storager) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	number := string(bodyBytes)

	num, err := strconv.Atoi(number)
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !luhn.Valid(num) {
		logger.WriteErrorLog("putOrder wrong format luhn number:" + number)
		rw.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	order, err := dbs.GetOrder(number)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	token := auth.GetTokenFromHeader(r)
	tokenData, err := dbs.GetAccessTokenData(token)
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if order.UserLogin != "" && order.UserLogin != tokenData.Login {
		logger.WriteErrorLog("putOrder order on other user")
		rw.WriteHeader(http.StatusConflict)
		return
	}

	if order.UserLogin != "" {
		logger.WriteInfoLog("putOrder order already on user")
		rw.WriteHeader(http.StatusOK)
		return
	}

	err = dbs.AddOrder(number, tokenData)
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.WriteInfoLog("putOrder order accepted on work")
	rw.WriteHeader(http.StatusAccepted)
}

// getOrders получение заказов
func getOrders(rw http.ResponseWriter, r *http.Request, dbs storage.Storager) {
	token := auth.GetTokenFromHeader(r)
	tokenData, err := dbs.GetAccessTokenData(token)
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	orders, err := dbs.GetOrders(tokenData.Login)
	if err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(orders) == 0 {
		logger.WriteErrorLog("putOrder GetOrders no rows returned")
		rw.WriteHeader(http.StatusNoContent)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	logger.WriteInfoLog("orders show")

	enc := json.NewEncoder(rw)
	if err = enc.Encode(orders); err != nil {
		logger.WriteErrorLog(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
