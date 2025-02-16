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
	"github.com/trysourcetool/sourcetool/backend/infra"
	"github.com/trysourcetool/sourcetool/backend/infra/mailer"
	"github.com/trysourcetool/sourcetool/backend/infra/signer"
	"github.com/trysourcetool/sourcetool/backend/infra/store"
	"github.com/trysourcetool/sourcetool/backend/logger"
	"github.com/trysourcetool/sourcetool/backend/postgres"
	"github.com/trysourcetool/sourcetool/backend/server"
	"github.com/trysourcetool/sourcetool/backend/ws"
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

	dep := infra.NewDependency(store.NewCE(db), signer.NewCE(), mailer.NewCE())

	if config.Config.Env == config.EnvLocal {
		if err := postgres.Migrate(); err != nil {
			logger.Logger.Fatal(err.Error())
		}

		if err := fixtures.Load(ctx, dep.Store); err != nil {
			logger.Logger.Fatal(err.Error())
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Logger.Info(fmt.Sprintf("Defaulting to port %s\n", port))
	}

	eg, egCtx := errgroup.WithContext(ctx)
	srv := &http.Server{
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      600 * time.Second,
		Handler:           server.New(dep),
		Addr:              fmt.Sprintf(":%s", port),
	}

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

		logger.Logger.Info("Closing WebSocket connections...")
		ws.GetConnManager().Close()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("server shutdown: %v", err)
		}

		logger.Logger.Info("Server shutdown complete")
		return nil
	})

	if err := eg.Wait(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Error(fmt.Sprintf("Error during shutdown: %v", err))
		os.Exit(1)
	}
}
