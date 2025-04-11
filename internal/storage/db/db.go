package db

import (
	"context"
	"time"

	"github.com/ramil063/firstgodiplom/internal/logger"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml"
)

type Storage struct{}

func Init(dbr dml.DataBaser) error {
	var err error

	if err = CheckPing(dbr); err != nil {
		logger.WriteErrorLog(err.Error())
		return err
	}

	err = CreateTables(dbr)
	return err
}

func CheckPing(dbr dml.DataBaser) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	return dbr.PingContext(ctx)
}

func CreateTables(dbr dml.DataBaser) error {
	var err error

	createTablesSQL := `
			--USERS
	CREATE TABLE IF NOT EXISTS public.users
	(
		id   		serial not null constraint pk_users_id primary key,
		name  		varchar(254),
		login    	varchar(100) not null constraint uidx_users_login unique,
		password 	varchar(100) not null,
		access_token            varchar(100),
		access_token_expired_at integer
	);
	COMMENT ON COLUMN public.users.id IS 'Идентификатор пользователя';
	COMMENT ON COLUMN public.users.name IS 'Имя пользователя';
	COMMENT ON COLUMN public.users.login IS 'Логин пользователя';
	COMMENT ON COLUMN public.users.password IS 'Пароль пользователя';
	COMMENT ON COLUMN public.users.access_token IS 'Токен пользователя';
	COMMENT ON COLUMN public.users.access_token_expired_at IS 'Токен пользователя истекает';

	        --ORDER
	CREATE TABLE IF NOT EXISTS public.order
	(
		id         	serial not null constraint pk_order_id primary key,
		number     	varchar(64) not null constraint uidx_order_number unique,
		accrual     double precision default 0 not null,
		status_id  	integer not null,
		created_at 	timestamp default current_timestamp not null,
    	uploaded_at varchar(64)                         not null,
    	user_id     integer not null constraint fk_order_user_id references users,
		check_accrual_after integer
	);
	COMMENT ON COLUMN public.order.id IS 'Идентификатор заказа';
	COMMENT ON COLUMN public.order.number IS 'Номер заказа';
	COMMENT ON COLUMN public.order.status_id IS 'Статус заказа';
	COMMENT ON COLUMN public.order.user_id IS 'Пользователь';
	COMMENT ON COLUMN public.order.created_at IS 'Дата создания записи';
	COMMENT ON COLUMN public.order.check_accrual_after IS 'Не проверять в сервисе акруал до этой даты';

	        --STATUS
	CREATE TABLE IF NOT EXISTS public.status
	(
		id    serial not null constraint pk_status_id primary key,
		alias varchar(100) not null,
		name  varchar(100) not null
	);
	COMMENT ON COLUMN public.status.id IS 'Идентификатор статуса';
	COMMENT ON COLUMN public.status.alias IS 'Псевдоним статуса';
	COMMENT ON COLUMN public.status.name IS 'Название статуса';
	TRUNCATE TABLE public.status;
	INSERT INTO public.status (id, alias, name) VALUES (2, 'PROCESSING', 'вознаграждение за заказ рассчитывается');
	INSERT INTO public.status (id, alias, name) VALUES (4, 'PROCESSED', 'данные по заказу проверены и информация о расчёте успешно получена');
	INSERT INTO public.status (id, alias, name) VALUES (1, 'NEW', 'заказ загружен в систему');
	INSERT INTO public.status (id, alias, name) VALUES (3, 'INVALID', 'система расчёта вознаграждений отказала в расчёте');

	        --WITHDRAW
	CREATE TABLE IF NOT EXISTS public.withdraw
	(
		id           serial not null constraint pk_withdraw_id primary key,
		sum          double precision default 0 not null,
		order_id     integer not null constraint fk_withdraw_order_id references public.order,
		processed_at varchar(100) not null,
		created_at   timestamp default current_timestamp not null
	);
	COMMENT ON COLUMN public.withdraw.id IS 'Идентификатор списания';
	COMMENT ON COLUMN public.withdraw.sum IS 'Сумма списания';
	COMMENT ON COLUMN public.withdraw.order_id IS 'Заказ';
	COMMENT ON COLUMN public.withdraw.processed_at IS 'Дата списания';
	COMMENT ON COLUMN public.withdraw.created_at IS 'Дата создания записи';

			--BALANCE
	CREATE TABLE IF NOT EXISTS public.balance
	(
		id      serial not null constraint balance_pk primary key,
		user_id integer not null constraint fk_balance_user_id references public.users,
		value   double precision default 0 not null
	);
	COMMENT ON COLUMN public.balance.id IS 'Идентификатор';
	COMMENT ON COLUMN public.balance.user_id IS 'Пользователь';
	COMMENT ON COLUMN public.balance.value IS 'Значение';
`
	_, err = dbr.ExecContext(context.Background(), createTablesSQL)
	return err
}
