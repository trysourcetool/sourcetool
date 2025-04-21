package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/trysourcetool/sourcetool/backend/cmd/internal"
	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
	"github.com/trysourcetool/sourcetool/backend/internal/mail"
	"github.com/trysourcetool/sourcetool/backend/internal/permission"
	"github.com/trysourcetool/sourcetool/backend/internal/postgres"
	"github.com/trysourcetool/sourcetool/backend/internal/pubsub"
	"github.com/trysourcetool/sourcetool/backend/internal/server"
	"github.com/trysourcetool/sourcetool/backend/internal/websocket"
)

func init() {
	config.Init()
	logger.Init()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	pqClient, err := postgres.Open()
	if err != nil {
		logger.Logger.Fatal("failed to open postgres", zap.Error(err))
	}

	redisClient, err := internal.OpenRedis()
	if err != nil {
		logger.Logger.Fatal("failed to open redis", zap.Error(err))
	}

	smtpClient, err := internal.OpenSMTP()
	if err != nil {
		logger.Logger.Fatal("failed to open smtp", zap.Error(err))
	}

	db := postgres.New(postgres.NewQueryLogger(pqClient))
	pubsub := pubsub.New(redisClient)
	mail := mail.New(smtpClient)
	wsManager := websocket.NewManager(ctx, db, pubsub)
	upgrader := internal.NewUpgrader()

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
	handler.Use(middleware.RequestID)
	handler.Use(middleware.Logger)
	handler.Use(middleware.Recoverer)
	handler.Use(middleware.Timeout(time.Duration(600) * time.Second))
	handler.Use(cors.New(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			if config.Config.IsCloudEdition {
				var pattern string

				switch config.Config.Env {
				case config.EnvProd:
					pattern = `^https://[a-zA-Z0-9-]+\.trysourcetool\.com$`
				case config.EnvStaging:
					pattern = `^https://[a-zA-Z0-9-]+\.staging\.trysourcetool\.com$`
				case config.EnvLocal:
					pattern = `^(http://[a-zA-Z0-9-]+\.local\.trysourcetool\.com:\d+|http://localhost:\d+)$`
				default:
					return false
				}

				matched, err := regexp.MatchString(pattern, origin)
				if err != nil {
					return false
				}

				return matched
			} else {
				// For self-hosted environments, we only need to check against the configured base URL
				normalizedOrigin := strings.TrimRight(origin, "/")
				normalizedBaseURL := strings.TrimRight(config.Config.BaseURL, "/")
				return normalizedOrigin == normalizedBaseURL
			}
		},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{},
		AllowCredentials: true,
		MaxAge:           0,
		Debug:            !(config.Config.Env == config.EnvProd),
	}).Handler)

	s := server.New(db, pubsub, mail, wsManager, permission.NewChecker(db), upgrader)
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

		if err := smtpClient.Close(); err != nil {
			logger.Logger.Sugar().Errorf("SMTP client close failed: %v", err)
		} else {
			logger.Logger.Sugar().Info("SMTP client gracefully stopped")
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
