package storage

import (
	accrualStorage "github.com/ramil063/firstgodiplom/cmd/gophermart/agent/accrual/storage"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/auth"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage/models/user/balance"
	"github.com/ramil063/firstgodiplom/internal/storage/db"
)

type Tokener interface {
	GetAccessTokenData(token string) (user.AccessTokenData, error)
	UpdateToken(u user.User, t auth.Token, expiredAt int64) error
}

type Userer interface {
	GetUser(login string) (user.User, error)
	AddUser(r user.Register, hash string) error
}

type Orderer interface {
	GetOrder(number string) (user.Order, error)
	GetOrders(login string) ([]user.Order, error)
	GetAllOrdersInStatuses(statuses []int) ([]user.OrderCheckAccrual, error)
	AddOrder(number string, tokenData user.AccessTokenData) error
	UpdateOrderAccrual(order accrualStorage.Order) error
	UpdateOrderCheckAccrualAfter(number string) error
}

type Balancer interface {
	GetBalance(login string) (balance.Balance, error)
	AddWithdraw(withdraw balance.Withdraw, login string) error
}

type Withdrawaler interface {
	GetWithdrawals(login string) ([]balance.Withdraw, error)
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
