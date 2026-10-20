package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/pur1/to_on/apps/api/auth/internal/config"
	"github.com/pur1/to_on/apps/api/auth/internal/database/mongo"
	"github.com/pur1/to_on/apps/api/auth/internal/handler"
	"github.com/pur1/to_on/apps/api/auth/internal/repository"
	"github.com/pur1/to_on/apps/api/auth/internal/server"
	"github.com/pur1/to_on/apps/api/auth/internal/service"
)

func main() {
	if err := run(); err != nil {
		slog.Error("api exited", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Local dev only: missing files are ignored and real env vars are never overridden.
	for _, f := range []string{"apps/api/.env", ".env"} {
		_ = godotenv.Load(f)
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	client, err := mongo.Connect(ctx, cfg.MongoURI)
	if err != nil {
		return err
	}
	defer client.Disconnect(context.Background())
	db := client.Database(cfg.MongoDB)

	users := repository.NewMongoUserRepository(db)
	tokens := repository.NewMongoTokenRepository(db)
	if err := users.EnsureIndexes(ctx); err != nil {
		return err
	}
	if err := tokens.EnsureIndexes(ctx); err != nil {
		return err
	}

	authSvc := service.NewAuthService(users, tokens, service.NewTokenIssuer(cfg.JWTSecret, cfg.AccessTTL), cfg.RefreshTTL)

	srv := server.New(":"+cfg.Port, server.Deps{
		Auth:        handler.NewAuthHandler(authSvc, cfg.CookieSecure),
		CORSOrigins: cfg.CORSOrigins,
	})

	errCh := make(chan error, 1)
	go func() {
		slog.Info("api listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
		close(errCh)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
