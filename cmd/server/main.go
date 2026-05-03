package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"tundraMarket/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"

	appstation "tundraMarket/internal/application/trading_station"
	sqlcdb "tundraMarket/internal/infrastructure/postgres/sqlc"
	stationinfrastructure "tundraMarket/internal/infrastructure/trading_station"
	httptransport "tundraMarket/internal/transport/http"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	pingCtx, cancel := context.WithTimeout(ctx, cfg.DBConnectTimeout)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		return err
	}

	queries := sqlcdb.New(pool)

	tradingStationRepo := stationinfrastructure.NewTradingStationRepo(queries)

	tradingStationUC := appstation.NewUseCase(tradingStationRepo)

	tradingStationHandler := httptransport.NewTradingStationHandler(tradingStationUC)

	handler := httptransport.NewRouter(httptransport.Dependencies{
		ReadinessCheck:        pool.Ping,
		TradingStationHandler: tradingStationHandler,
	})

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           handler,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		ReadTimeout:       cfg.ReadTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("starting server on %s", cfg.HTTPAddr)
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}
}
