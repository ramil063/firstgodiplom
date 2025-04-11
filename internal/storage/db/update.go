package db

import (
	"errors"
	orderConstants "github.com/ramil063/firstgodiplom/internal/constants/order"
	"github.com/ramil063/firstgodiplom/internal/env"
	"time"

	accrualStorage "github.com/ramil063/firstgodiplom/cmd/gophermart/agent/accrual/storage"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	statusConstants "github.com/ramil063/firstgodiplom/internal/constants/status"
	"github.com/ramil063/firstgodiplom/internal/logger"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml"
)

func (s *Storage) UpdateToken(u user.User, t auth.Token, expiredAt int64) error {

	result, err := dml.UpdateToken(&dml.DBRepository, u.Login, t.Token, expiredAt)

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
	result, err := dml.UpdateOrderAccrual(&dml.DBRepository, order.Number, order.Accrual, internalOrderStatus)

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
