package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Unlites/wishlist/internal/common/parser"
	"github.com/Unlites/wishlist/internal/config"
	"github.com/Unlites/wishlist/internal/handlers/http/middleware"
	userHandlerHttp "github.com/Unlites/wishlist/internal/handlers/http/user"
	wishHandlerHttp "github.com/Unlites/wishlist/internal/handlers/http/wish"
	"github.com/Unlites/wishlist/internal/infra/hasher"
	userRepoPg "github.com/Unlites/wishlist/internal/infra/repositories/user/postgres"
	wishRepoPg "github.com/Unlites/wishlist/internal/infra/repositories/wish/postgres"
	"github.com/Unlites/wishlist/internal/infra/tokenmanager"
	"github.com/Unlites/wishlist/internal/services/user"
	"github.com/Unlites/wishlist/internal/services/wish"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.New()

	ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: parser.ParseSlogLevel(cfg.LogLevel),
	}))

	slog.SetDefault(logger)

	pgPool, err := pgxpool.New(ctx, cfg.Postgres.DSN)
	if err != nil {
		slog.Error("failed to create pg pool", "detail", err)
		os.Exit(1)
	}

	if err := pgPool.Ping(ctx); err != nil {
		slog.Error("failed to ping pg", "detail", err)
		os.Exit(1)
	}

	userRepo := userRepoPg.NewUserRepositoryPostgres(pgPool)
	wishRepo := wishRepoPg.NewWishRepositoryPostgres(pgPool)

	hasher := hasher.NewBcryptHasher(cfg.HashCost)
	tokenManager := tokenmanager.NewJWTTokenManager(cfg.JWT.SecretKey, cfg.JWT.ExpirationTime)

	userService := user.NewUserService(userRepo, hasher, tokenManager)
	wishService := wish.NewWishService(wishRepo)

	middlewareProvider := middleware.NewMiddlewareProvider(tokenManager)

	mux := http.NewServeMux()

	wishHandler := wishHandlerHttp.NewWishHandler(wishService, middlewareProvider)
	userHandler := userHandlerHttp.NewUserHandler(userService)

	wishHandler.RegisterRoutes(mux, "/api/v1/users/{user_id}/wishes")
	userHandler.RegisterRoutes(mux, "/api/v1/users")

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      mux,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("srv.ListenAndServe", "detail", err)
			os.Exit(1)
		}
	}()

	slog.Warn("service started", "address", srv.Addr)

	notifyCtx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	<-notifyCtx.Done()
	slog.Warn("service gracefully stopping...")

	shutDownCtx, cancel := context.WithTimeout(ctx, cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutDownCtx); err != nil {
		slog.Error("failed to shutdown server", "detail", err)
		os.Exit(1)
	}

	slog.Warn("service stopped")
}
