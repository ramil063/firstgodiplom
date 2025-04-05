package db

import (
	"errors"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
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
