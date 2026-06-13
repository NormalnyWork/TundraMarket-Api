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

	appauth "tundraMarket/internal/application/auth"
	apporder "tundraMarket/internal/application/order"
	appproduct "tundraMarket/internal/application/product"
	appstation "tundraMarket/internal/application/trading_station"
	admininfrastructure "tundraMarket/internal/infrastructure/admin"
	authinfrastructure "tundraMarket/internal/infrastructure/auth"
	nomadinfrastructure "tundraMarket/internal/infrastructure/nomad"
	orderinfrastructure "tundraMarket/internal/infrastructure/order"
	productinfrastructure "tundraMarket/internal/infrastructure/product"
	stationinfrastructure "tundraMarket/internal/infrastructure/trading_station"
	httptransport "tundraMarket/internal/transport/http"

	sqlcdb "tundraMarket/internal/infrastructure/postgres/sqlc"
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

	nomadRepo := nomadinfrastructure.NewNomadRepo(queries)
	productRepo := productinfrastructure.NewProductRepo(queries)
	tradingStationRepo := stationinfrastructure.NewTradingStationRepo(queries)
	orderRepo := orderinfrastructure.NewOrderRepo(pool, queries)
	adminRepo := admininfrastructure.NewAdminRepo(queries)

	tokenIssuer := authinfrastructure.NewTokenIssuer(cfg.AuthTokenSecret, cfg.AuthTokenTTL)
	passwordVerifier := authinfrastructure.NewPasswordVerifier()

	authUC := appauth.NewUseCase(
		nomadRepo,
		tradingStationRepo,
		adminRepo,
		tokenIssuer,
		passwordVerifier,
	)
	productUC := appproduct.NewUseCase(productRepo)
	tradingStationUC := appstation.NewUseCase(tradingStationRepo)
	orderUC := apporder.NewUseCase(orderRepo, nomadRepo, tradingStationRepo, productRepo)

	authHandler := httptransport.NewAuthHandler(authUC)
	productHandler := httptransport.NewProductHandler(productUC)
	tradingStationHandler := httptransport.NewTradingStationHandler(tradingStationUC)
	orderHandler := httptransport.NewOrderHandler(orderUC)

	handler := httptransport.NewRouter(httptransport.Dependencies{
		ReadinessCheck:        pool.Ping,
		TokenIssuer:           tokenIssuer,
		AuthHandler:           authHandler,
		ProductHandler:        productHandler,
		OrderHandler:          orderHandler,
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
