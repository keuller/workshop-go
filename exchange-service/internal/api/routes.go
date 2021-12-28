package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimid "github.com/go-chi/chi/v5/middleware"

	"github.com/keuller/exchange/internal/controller"
)

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
	v1 := chi.NewRouter()
	v1.Get("/exchange/currencies", controller.GetCurrencies)
	v1.Get("/exchange/quotation", controller.GetQuotation)
	return v1
}
