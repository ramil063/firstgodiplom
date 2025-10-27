package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/middlewares"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/handlers/router/balance"
	"github.com/ramil063/firstgodiplom/cmd/gophermart/server/storage"
	"github.com/ramil063/firstgodiplom/internal/logger"
)

func Router(s storage.Storager) chi.Router {
	r := chi.NewRouter()

	r.Use(logger.ResponseLogger)
	r.Use(logger.RequestLogger)
	r.Use(middlewares.GZIPMiddleware)

	r.Route("/api/user/register", func(r chi.Router) {
		r.Use(middlewares.CheckContentTypeMiddleware("application/json"))

		userRegisterHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
			userRegister(rw, r, s)
		}
		r.With(middlewares.CheckPostMethodMw).Post("/", userRegisterHandlerFunction)
	})

	r.Route("/api/user/login", func(r chi.Router) {
		r.Use(middlewares.CheckContentTypeMiddleware("application/json"))

		userLoginHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
			userLogin(rw, r, s)
		}
		r.With(middlewares.CheckPostMethodMw).Post("/", userLoginHandlerFunction)
	})

	r.Route("/api/user/orders", func(r chi.Router) {
		r.Use(middlewares.CheckAuthMiddleware(s))

		putOrderHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
			putOrder(rw, r, s)
		}
		r.With(middlewares.CheckPostMethodMw).
			With(middlewares.CheckContentTypeMiddleware("text/plain")).
			Post("/", putOrderHandlerFunction)

		getOrderHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
			getOrders(rw, r, s)
		}
		r.With(middlewares.CheckGetMethodMw).Get("/", getOrderHandlerFunction)
	})

	r.Route("/api/user/balance", func(r chi.Router) {
		r.Use(middlewares.CheckAuthMiddleware(s))

		getBalanceHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
			balance.GetBalance(rw, r, s)
		}
		r.With(middlewares.CheckGetMethodMw).Get("/", getBalanceHandlerFunction)

		r.Route("/withdraw", func(r chi.Router) {
			addWithdrawHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
				balance.AddWithdraw(rw, r, s)
			}
			r.With(middlewares.CheckPostMethodMw).
				With(middlewares.CheckContentTypeMiddleware("application/json")).
				Post("/", addWithdrawHandlerFunction)
		})
	})

	r.Route("/api/user/withdrawals", func(r chi.Router) {
		r.Use(middlewares.CheckAuthMiddleware(s))

		getWithdrawHandlerFunction := func(rw http.ResponseWriter, r *http.Request) {
			getWithdrawals(rw, r, s)
		}
		r.With(middlewares.CheckGetMethodMw).Get("/", getWithdrawHandlerFunction)
	})
	return r
}
