package main

import (
	"net/http"

	accrualHandler "github.com/ramil063/firstgodiplom/cmd/gophermart/agent/accrual/handlers"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/flags"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/router"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	"github.com/ramil063/firstgodiplom/internal/logger"
	"github.com/ramil063/firstgodiplom/internal/storage/db"
	"github.com/ramil063/firstgodiplom/internal/storage/db/dml"
)

func main() {

	var err error
	if err = logger.Initialize(); err != nil {
		panic(err)
	}
	flags.ParseFlags()

	logger.WriteInfoLog("--------------START SERVER-------------")

	if flags.DatabaseURI != "" {
		rep, err := dml.NewRepository()
		if err != nil {
			logger.WriteErrorLog(err.Error())
			return
		}
		dml.DBRepository = *rep
		err = db.Init(&dml.DBRepository)
		defer dml.DBRepository.Pool.Close()
		if err != nil {
			logger.WriteErrorLog(err.Error())
			return
		}
	}

	s := storage.NewDBStorage()
	c := accrualHandler.NewClient()
	accrualHandler.OrdersProcess(c, s)

	if err = http.ListenAndServe(flags.RunAddress, router.Router(s)); err != nil {
		panic(err)
	}
}
