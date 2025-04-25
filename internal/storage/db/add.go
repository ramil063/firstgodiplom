package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
	orderStatus "github.com/ramil063/firstgodiplom/internal/constants/status"
	internalErrors "github.com/ramil063/firstgodiplom/internal/errors"
	"github.com/ramil063/firstgodiplom/internal/logger"
	balanceRepository "github.com/ramil063/firstgodiplom/internal/storage/db/dml/balance"
	orderRepository "github.com/ramil063/firstgodiplom/internal/storage/db/dml/order"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository"
	userRepository "github.com/ramil063/firstgodiplom/internal/storage/db/dml/user"
	withdrawRepository "github.com/ramil063/firstgodiplom/internal/storage/db/dml/withdraw"
)

// AddUserData добавить пользователя и добавить его баланс
func (s *Storage) AddUserData(ctx context.Context, register user.Register, passwordHash string) error {

	tx, err := repository.DBRepository.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	result, err := userRepository.AddUser(ctx, tx, register.Login, passwordHash, register.Name)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if result == nil {
		_ = tx.Rollback(ctx)
		logger.WriteErrorLog("AddUser error in sql empty result")
		return errors.New("error in sql empty result")
	}

	rows := result.RowsAffected()
	if rows != 1 {
		_ = tx.Rollback(ctx)
		logger.WriteErrorLog("AddUser error expected to affect 1 row")
		return errors.New("expected to affect 1 row")
	}

	result, err = balanceRepository.AddBalance(ctx, tx, register.Login)
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if result == nil {
		_ = tx.Rollback(ctx)
		logger.WriteErrorLog("AddBalance error in sql empty result")
		return errors.New("error in sql empty result")
	}

	rows = result.RowsAffected()
	if rows != 1 {
		_ = tx.Rollback(ctx)
		logger.WriteErrorLog("AddBalance error expected to affect 1 row")
		return errors.New("expected to affect 1 row")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

// AddOrder добавить заказ
func (s *Storage) AddOrder(ctx context.Context, number string, tokenData user.AccessTokenData) error {

	u, err := s.GetUser(ctx, tokenData.Login)
	if err != nil {
		return err
	}

	now := time.Now()
	rfc3339Time := now.Format(time.RFC3339)

	result, err := orderRepository.AddOrder(ctx, &repository.DBRepository, number, 0, orderStatus.NewID, rfc3339Time, u.ID)
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

// AddWithdrawFromBalance списание баллов с баланса
func (s *Storage) AddWithdrawFromBalance(ctx context.Context, withdraw balance.Withdraw, login string) error {

	now := time.Now()
	rfc3339Time := now.Format(time.RFC3339)

	tx, err := repository.DBRepository.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	balanceData, err := balanceRepository.GetBalanceForUpdate(ctx, tx, login)
	if err != nil {
		_ = tx.Rollback(ctx)
		logger.WriteErrorLog("GetBalanceForUpdate error in sql empty result")
		return err
	}

	if balanceData.Current < 0 {
		_ = tx.Rollback(ctx)
		return errors.New("balance under 0")
	}

	if balanceData.Current < withdraw.Sum {
		_ = tx.Rollback(ctx)
		return internalErrors.ErrNotEnoughBalance
	}

	result, err := withdrawRepository.AddWithdraw(ctx, tx, withdraw.OrderNumber, withdraw.Sum, rfc3339Time, login)
	if err != nil {
		errTx := tx.Rollback(ctx)
		if errTx != nil {
			return errTx
		}
		return err
	}

	if result == nil {
		_ = tx.Rollback(ctx)
		logger.WriteErrorLog("AddWithdraw error in sql empty result")
		return errors.New("error in sql empty result")
	}

	rows := result.RowsAffected()
	if rows != 1 {
		_ = tx.Rollback(ctx)
		logger.WriteErrorLog("AddWithdraw RowsAffected error expected to affect 1 row")
		return errors.New("expected to affect 1 row")
	}

	_, err = balanceRepository.OperatingBalance(ctx, tx, withdraw.Sum, "-", login)
	if err != nil {
		errTx := tx.Rollback(ctx)
		if errTx != nil {
			return errTx
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
