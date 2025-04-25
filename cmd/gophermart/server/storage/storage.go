package storage

import (
	"context"

	accrualStorage "github.com/ramil063/firstgodiplom/cmd/gophermart/agent/accrual/storage"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
	"github.com/ramil063/firstgodiplom/internal/storage/db"
)

type Tokener interface {
	GetAccessTokenData(ctx context.Context, token string) (user.AccessTokenData, error)
	UpdateToken(ctx context.Context, login string, t auth.Token, expiredAt int64) error
}

type Userer interface {
	GetUser(ctx context.Context, login string) (user.User, error)
	AddUserData(ctx context.Context, r user.Register, hash string) error
}

type Orderer interface {
	GetOrder(ctx context.Context, number string) (user.Order, error)
	GetOrders(ctx context.Context, login string) ([]user.Order, error)
	GetAllOrdersInStatuses(ctx context.Context, statuses []int) ([]user.OrderCheckAccrual, error)
	AddOrder(ctx context.Context, number string, tokenData user.AccessTokenData) error
	UpdateOrderAccrual(ctx context.Context, order accrualStorage.Order) error
	UpdateOrderCheckAccrualAfter(ctx context.Context, number string) error
}

type Balancer interface {
	GetBalance(ctx context.Context, login string) (balance.Balance, error)
	AddWithdrawFromBalance(ctx context.Context, withdraw balance.Withdraw, login string) error
}

type Withdrawaler interface {
	GetWithdrawals(ctx context.Context, login string) ([]balance.Withdraw, error)
}

type Storager interface {
	Tokener
	Userer
	Orderer
	Balancer
	Withdrawaler
}

func NewDBStorage() Storager {
	return &db.Storage{}
}
