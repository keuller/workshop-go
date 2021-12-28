package api

import (
	"net/http"

	"github.com/keuller/exchange/internal/controller"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
)

func configureRoutes() *bunrouter.CompatRouter {
	router := bunrouter.New(
		bunrouter.WithMiddleware(reqlog.NewMiddleware()),
	).Compat()

	router.GET("/liveness", heartbeat)
	router.WithGroup("/v1", v1Routes)
	// setMiddlewares(router)
	// router.Mount("/v1", v1Routes())
	return router
}

/**
func setMiddlewares(mux *chi.Mux) {
	// mux.Use(chimid.StripSlashes)
	// mux.Use(chimid.Logger)
	// mux.Use(chimid.Heartbeat("/liveness"))
	mux.Use(chimid.RequestID)
	mux.Use(chimid.Timeout(30 * time.Second))
}
*/

func v1Routes(group *bunrouter.CompatGroup) {
	group.GET("/exchange/currencies", controller.GetCurrencies)
	group.GET("/exchange/quotation", controller.GetQuotation)
}

func heartbeat(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("UP"))
}
