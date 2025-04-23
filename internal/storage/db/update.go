package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"

	accrualStorage "github.com/ramil063/firstgodiplom/cmd/gophermart/agent/accrual/storage"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/auth"
	orderConstants "github.com/ramil063/firstgodiplom/internal/constants/order"
	statusConstants "github.com/ramil063/firstgodiplom/internal/constants/status"
	"github.com/ramil063/firstgodiplom/internal/env"
	"github.com/ramil063/firstgodiplom/internal/logger"
	balanceRepository "github.com/ramil063/firstgodiplom/internal/storage/db/dml/balance"
	orderRepository "github.com/ramil063/firstgodiplom/internal/storage/db/dml/order"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository"
	userRepository "github.com/ramil063/firstgodiplom/internal/storage/db/dml/user"
)

// UpdateToken обновить токен авторизации
func (s *Storage) UpdateToken(login string, t auth.Token, expiredAt int64) error {

	result, err := userRepository.UpdateToken(&repository.DBRepository, login, t.Token, expiredAt)

	if err != nil {
		return err
	}
	if result == nil {
		logger.WriteErrorLog("error in sql empty result")
		return errors.New("error in sql empty result")
	}

	rows := result.RowsAffected()
	if rows != 1 {
		logger.WriteErrorLog("error expected to affect 1 row")
		return errors.New("expected to affect 1 row")
	}
	return nil
}

// UpdateOrderAccrual обновить данные по начислению данных в заказе
func (s *Storage) UpdateOrderAccrual(orderFromAccrual accrualStorage.Order) error {

	internalOrderStatus := statusConstants.AccrualStatusesOrderStatusesMap[orderFromAccrual.Status]

	tx, err := repository.DBRepository.Pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	resultUpdateOrderAccrual, err := orderRepository.UpdateOrderAccrual(tx, orderFromAccrual.Order, orderFromAccrual.Accrual, internalOrderStatus)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return err
	}
	if resultUpdateOrderAccrual == nil {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("UpdateOrderAccrual error in sql empty result")
		return errors.New("error in sql empty result")
	}
	rows := resultUpdateOrderAccrual.RowsAffected()
	if rows != 1 {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("UpdateOrderAccrual error expected to affect 1 row")
		return errors.New("expected to affect 1 row")
	}

	orderEntity, err := orderRepository.GetOrder(tx, orderFromAccrual.Order)
	if err != nil {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("GetOrder error in sql")
		return err
	}

	balance, err := balanceRepository.GetBalance(tx, orderEntity.UserLogin)
	if err != nil {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("GetBalance error in sql")
		return err
	}
	if balance.ID == 0 {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("GetBalance error in sql")
		return errors.New("error in sql empty result")
	}

	_, err = balanceRepository.OperatingBalance(tx, orderFromAccrual.Accrual, "+", orderEntity.UserLogin)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return err
	}

	err = tx.Commit(context.Background())
	return err
}

// UpdateOrderCheckAccrualAfter обновить дату повторного запроса данных заказа в акруал
func (s *Storage) UpdateOrderCheckAccrualAfter(orderNumber string) error {

	checkAfterUnix := time.Now().Unix()
	if env.AppEnv == "PROD" {
		checkAfterUnix += orderConstants.CheckAccrualAfterSeconds
	}
	result, err := orderRepository.UpdateOrderCheckAccrualAfter(&repository.DBRepository, orderNumber, checkAfterUnix)

	if err != nil {
		return err
	}
	if result == nil {
		logger.WriteErrorLog("UpdateOrderCheckAccrualAfter error in sql empty result")
		return errors.New("error in sql empty result")
	}

	rows := result.RowsAffected()
	if rows != 1 {
		logger.WriteErrorLog("UpdateOrderCheckAccrualAfter error expected to affect 1 row")
		return errors.New("expected to affect 1 row")
	}
	return nil
}
