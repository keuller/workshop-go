package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"github.com/keuller/account/internal/application"
	"github.com/keuller/account/internal/common"
	"github.com/keuller/account/internal/controller"
	"github.com/keuller/account/internal/infra"
	"github.com/keuller/account/internal/infra/repository"
)

type operation func(http.ResponseWriter, *http.Request) error

func handleError(fn operation) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		if err := fn(res, req); err != nil {
			if errors.Is(err, common.ErrNotFound) {
				http.Error(res, err.Error(), http.StatusNotFound)
			} else if errors.Is(err, common.ErrValidation) || errors.Is(err, validator.ValidationErrors{}) {
				http.Error(res, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(res, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}

func configureRoutes() *chi.Mux {
	router := chi.NewRouter()
	setMiddlewares(router)

	accountRepo := repository.NewAccountRepository(infra.Connection())
	exchangeRepo := repository.NewExchangeRepository()
	accountSvc := application.NewAccountService(accountRepo, exchangeRepo)
	accountCtrl := controller.NewAccountController(accountSvc)

	router.Post("/v1/accounts", handleError(accountCtrl.CreateAccountHandler))
	router.Get("/v1/accounts/{code}/balance", accountCtrl.GetBalanceHandler)
	router.Post("/v1/accounts/deposit", accountCtrl.Deposit)
	router.Patch("/v1/accounts/transfer", accountCtrl.Transfer)
	return router
}

func setMiddlewares(mux *chi.Mux) {
	mux.Use(chimid.StripSlashes)
	mux.Use(chimid.Logger)
	mux.Use(chimid.Heartbeat("/liveness"))
	mux.Use(chimid.RequestID)
	mux.Use(chimid.Timeout(30 * time.Second))
}
