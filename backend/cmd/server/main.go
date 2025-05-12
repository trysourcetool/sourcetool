package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/trysourcetool/sourcetool/backend/cmd/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/encrypt"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
	"github.com/trysourcetool/sourcetool/backend/internal/permission"
	"github.com/trysourcetool/sourcetool/backend/internal/postgres"
	"github.com/trysourcetool/sourcetool/backend/internal/pubsub"
	"github.com/trysourcetool/sourcetool/backend/internal/server"
	"github.com/trysourcetool/sourcetool/backend/internal/ws"
)

func init() {
	config.Init()
	logger.Init()
}

func main() {
	autoMigrateFlag := flag.Bool("auto-migrate", false, "run migrations before starting the server")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if *autoMigrateFlag {
		if err := postgres.Migrate("migrations"); err != nil {
			logger.Logger.Fatal("auto migration failed", zap.Error(err))
		}
	}

	if err := internal.CheckLicense(); err != nil {
		logger.Logger.Fatal("license check failed", zap.Error(err))
	}

	pqClient, err := postgres.Open()
	if err != nil {
		logger.Logger.Fatal("failed to open postgres", zap.Error(err))
	}

	redisClient, err := internal.OpenRedis()
	if err != nil {
		logger.Logger.Fatal("failed to open redis", zap.Error(err))
	}

	db := postgres.New(pqClient)
	pubsub := pubsub.New(redisClient)
	wsManager := ws.NewManager(ctx, db, pubsub)
	upgrader := internal.NewUpgrader()
	encryptor, err := encrypt.NewEncryptor()
	if err != nil {
		logger.Logger.Fatal("failed to create encryptor", zap.Error(err))
	}

	if config.Config.Env == config.EnvLocal {
		if err := internal.LoadFixtures(ctx, db); err != nil {
			logger.Logger.Fatal(err.Error())
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Logger.Info(fmt.Sprintf("Defaulting to port %s\n", port))
	}

	handler := chi.NewRouter()
	s := server.New(db, pubsub, wsManager, permission.NewChecker(db), upgrader, encryptor)
	s.Install(handler)

	srv := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      600 * time.Second,
		Handler:           handler,
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
			logger.Logger.Sugar().Errorf("WebSocket manager graceful shutdown failed: %v", err)
		} else {
			logger.Logger.Sugar().Info("WebSocket manager gracefully stopped")
		}

		if err := redisClient.Close(); err != nil {
			logger.Logger.Sugar().Errorf("Redis client close failed: %v", err)
		} else {
			logger.Logger.Sugar().Info("Redis client gracefully stopped")
		}

		if err := pqClient.Close(); err != nil {
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
