package httptransport

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Dependencies struct {
	ReadinessCheck        func(context.Context) error
	AuthHandler           *AuthHandler
	ProductHandler        *ProductHandler
	TradingStationHandler *TradingStationHandler
}

func NewRouter(deps Dependencies) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", handleHealth)
	r.Get("/readyz", handleReadiness(deps.ReadinessCheck))

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/user/auth", deps.AuthHandler.Auth)
		r.Get("/user/catalog", deps.ProductHandler.Catalog)
		r.Get("/trading-stations/list", deps.TradingStationHandler.List)
	})

	return r
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleReadiness(check func(context.Context) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if check == nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "readiness_check_not_configured"})
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		if err := check(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not_ready"})
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
	}
}
