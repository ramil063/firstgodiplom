package user

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	internalErrors "github.com/ramil063/firstgodiplom/internal/errors"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml/repository"
)

// AddUser добавить пользователя
func AddUser(ctx context.Context, tx pgx.Tx, login string, password string, name string) (pgconn.CommandTag, error) {
	exec, err := tx.Exec(
		ctx,
		"INSERT INTO users (login, password, name) VALUES ($1, $2, $3)",
		login,
		password,
		name)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, internalErrors.ErrUniqueViolation
		}
		return nil, err
	}
	return exec, nil
}

// UpdateToken обновить токен авторизации пользователя
func UpdateToken(ctx context.Context, dbr *repository.Repository, login string, token string, expiredAt int64) (pgconn.CommandTag, error) {
	exec, err := dbr.ExecContext(
		ctx,
		"UPDATE users SET access_token = $1, access_token_expired_at = $2 WHERE login = $3",
		token,
		expiredAt,
		login)

	if err != nil {
		return nil, internalErrors.NewDBError(err)
	}
	return exec, nil
}
