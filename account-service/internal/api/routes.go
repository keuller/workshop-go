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
	"github.com/keuller/account/internal/domain"
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
	router.Mount("/v1", v1Routes())
	return router
}

func setMiddlewares(mux *chi.Mux) {
	mux.Use(chimid.StripSlashes)
	mux.Use(chimid.Logger)
	mux.Use(chimid.Heartbeat("/liveness"))
	mux.Use(chimid.RequestID)
	mux.Use(chimid.Timeout(30 * time.Second))
}

func v1Routes() http.Handler {
	accountRepo := repository.NewAccountRepository(infra.Connection())
	exchangeRepo := repository.NewExchangeRepository()
	accountSvc := domain.NewAccountService(accountRepo, exchangeRepo)
	accountCtrl := application.NewAccountController(accountSvc)

	v1 := chi.NewRouter()

	v1.Post("/accounts", handleError(accountCtrl.CreateAccountHandler))
	v1.Get("/accounts/{code}/balance", accountCtrl.GetBalanceHandler)
	v1.Post("/accounts/deposit", accountCtrl.Deposit)
	v1.Patch("/accounts/transfer", accountCtrl.Transfer)

	return v1
}
