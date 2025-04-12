package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
	orderStatus "github.com/ramil063/firstgodiplom/internal/constants/status"
	"github.com/ramil063/firstgodiplom/internal/logger"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml"
)

func (s *Storage) AddUserData(register user.Register, passwordHash string) error {

	tx, err := dml.DBRepository.Pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}

	result, err := dml.AddUser(tx, register.Login, passwordHash, register.Name)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return err
	}
	if result == nil {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("AddUser error in sql empty result")
		return errors.New("error in sql empty result")
	}

	rows := result.RowsAffected()
	if rows != 1 {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("AddUser error expected to affect 1 row")
		return errors.New("expected to affect 1 row")
	}

	result, err = dml.AddBalance(tx, register.Login)
	if err != nil {
		_ = tx.Rollback(context.Background())
		return err
	}
	if result == nil {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("AddBalance error in sql empty result")
		return errors.New("error in sql empty result")
	}

	rows = result.RowsAffected()
	if rows != 1 {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("AddBalance error expected to affect 1 row")
		return errors.New("expected to affect 1 row")
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddOrder(number string, tokenData user.AccessTokenData) error {

	u, err := s.GetUser(tokenData.Login)
	if err != nil {
		return err
	}

	now := time.Now()
	rfc3339Time := now.Format(time.RFC3339)

	result, err := dml.AddOrder(&dml.DBRepository, number, 0, orderStatus.NewID, rfc3339Time, u.ID)
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

func (s *Storage) AddWithdraw(withdraw balance.Withdraw, login string) error {

	now := time.Now()
	rfc3339Time := now.Format(time.RFC3339)

	tx, err := dml.DBRepository.Pool.BeginTx(context.Background(), pgx.TxOptions{})
	if err != nil {
		return err
	}

	result, err := dml.AddWithdraw(tx, withdraw.OrderNumber, withdraw.Sum, rfc3339Time, login)
	if err != nil {
		errTx := tx.Rollback(context.Background())
		if errTx != nil {
			return errTx
		}
		return err
	}

	if result == nil {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("error in sql empty result")
		return errors.New("error in sql empty result")
	}

	rows := result.RowsAffected()
	if rows != 1 {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("error expected to affect 1 row")
		return errors.New("expected to affect 1 row")
	}

	result, err = dml.OperatingBalance(tx, withdraw.Sum, "+", login)
	if err != nil {
		errTx := tx.Rollback(context.Background())
		if errTx != nil {
			return errTx
		}
		return err
	}

	if result == nil {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("error in sql empty result")
		return errors.New("error in sql empty result")
	}

	rows = result.RowsAffected()
	if rows != 1 {
		_ = tx.Rollback(context.Background())
		logger.WriteErrorLog("error expected to affect 1 row")
		return errors.New("expected to affect 1 row")
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}
