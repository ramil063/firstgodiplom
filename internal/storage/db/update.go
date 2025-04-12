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
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml"
)

func (s *Storage) UpdateToken(login string, t auth.Token, expiredAt int64) error {

	result, err := dml.UpdateToken(&dml.DBRepository, login, t.Token, expiredAt)

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

func (s *Storage) UpdateOrderAccrual(orderFromAccrual accrualStorage.Order) error {

	order, err := s.GetOrder(orderFromAccrual.Order)
	if err != nil {
		return err
	}
	internalOrderStatus := statusConstants.AccrualStatusesOrderStatusesMap[orderFromAccrual.Status]

	tx, err := dml.DBRepository.Pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}

	resultUpdateOrderAccrual, err := dml.UpdateOrderAccrual(tx, order.Number, orderFromAccrual.Accrual, internalOrderStatus)
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

	orderEntity, err := dml.GetOrder(tx, order.Number)
	if err != nil {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("GetOrder error in sql")
		return err
	}

	balance, err := dml.GetBalance(tx, orderEntity.UserLogin)
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

	resultOperatingBalance, err := dml.OperatingBalance(tx, orderFromAccrual.Accrual, "+", orderEntity.UserLogin)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return err
	}
	if resultOperatingBalance == nil {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("OperatingBalance error in sql empty result")
		return errors.New("error in sql empty result")
	}

	rows = resultOperatingBalance.RowsAffected()
	if rows != 1 {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("OperatingBalance error expected to affect 1 row")
		return errors.New("expected to affect 1 row")
	}
	err = tx.Commit(context.Background())
	return err
}

func (s *Storage) UpdateOrderCheckAccrualAfter(orderNumber string) error {

	checkAfterUnix := time.Now().Unix()
	if env.AppEnv == "PROD" {
		checkAfterUnix += orderConstants.CheckAccrualAfterSeconds
	}
	result, err := dml.UpdateOrderCheckAccrualAfter(&dml.DBRepository, orderNumber, checkAfterUnix)

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
