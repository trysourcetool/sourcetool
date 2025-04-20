package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/trysourcetool/sourcetool/backend/config"
	_ "github.com/trysourcetool/sourcetool/backend/docs"
	"github.com/trysourcetool/sourcetool/backend/fixtures"
	"github.com/trysourcetool/sourcetool/backend/internal/app/port"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/db/postgres"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/email/smtp"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/pubsub/redis"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/ws/manager"
	"github.com/trysourcetool/sourcetool/backend/internal/transport"
	"github.com/trysourcetool/sourcetool/backend/logger"
)

func init() {
	config.Init()
	logger.Init()
}

// @title Sourcetool API
// @version 1.0
// @description Sourcetool's API documentation
// @termsOfService http://swagger.io/terms/
// @host https://api.trysourcetool.com
// @BasePath /api/v1.
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, err := postgres.New()
	if err != nil {
		logger.Logger.Fatal("failed to open postgres", zap.Error(err))
	}

	logger.Logger.Sugar().Infof("Connected to Postgres at %s:%s", config.Config.Postgres.Host, config.Config.Postgres.Port)

	redisClient, err := redis.NewClientCE()
	if err != nil {
		logger.Logger.Fatal("failed to open redis", zap.Error(err))
	}

	logger.Logger.Sugar().Infof("Connected to Redis at %s:%s", config.Config.Redis.Host, config.Config.Redis.Port)

	repo := postgres.NewRepositoryCE(db)
	smtpMailer := smtp.NewMailerCE()
	wsManager := manager.NewManager(ctx, repo, redisClient)

	// Initialize infra dependency container
	dep := port.NewDependencies(repo, smtpMailer, redisClient, wsManager)

	if config.Config.Env == config.EnvLocal {
		if err := fixtures.Load(ctx, dep.Repository); err != nil {
			logger.Logger.Fatal(err.Error())
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Logger.Info(fmt.Sprintf("Defaulting to port %s\n", port))
	}

	srv := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      600 * time.Second,
		Handler:           transport.NewRouter(dep).Build(),
		Addr:              fmt.Sprintf(":%s", port),
	}

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		logger.Logger.Info(fmt.Sprintf("Listening on port %s\n", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("HTTP server error: %v", err)
		}
		return nil
	})
	eg.Go(func() error {
		<-egCtx.Done()
		logger.Logger.Info("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		// Attempt to gracefully shut down the server first.
		var shutdownErr error
		if err := srv.Shutdown(shutdownCtx); err != nil {
			logger.Logger.Error("Server shutdown error", zap.Error(err))
			shutdownErr = fmt.Errorf("server shutdown: %v", err)
		}

		if err := wsManager.Close(); err != nil {
			logger.Logger.Sugar().Errorf("WebSocket Manager graceful shutdown failed: %v", err)
		} else {
			logger.Logger.Sugar().Info("WebSocket Manager gracefully stopped")
		}

		if err := redisClient.Close(); err != nil {
			logger.Logger.Sugar().Errorf("Redis client close failed: %v", err)
		} else {
			logger.Logger.Sugar().Info("Redis client gracefully stopped")
		}

		if err := db.Close(); err != nil {
			logger.Logger.Sugar().Errorf("DB connection close failed: %v", err)
		} else {
			logger.Logger.Sugar().Info("DB connection gracefully stopped")
		}

		logger.Logger.Info("Server shutdown complete")
		// Return the server shutdown error if it happened.
		return shutdownErr
	})

	if err := eg.Wait(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Error(fmt.Sprintf("Error during shutdown: %v", err))
		os.Exit(1)
	}
}
