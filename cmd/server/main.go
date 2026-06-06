package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ssjoy/website-health-checker/internal/config"
	"github.com/ssjoy/website-health-checker/internal/database"
	"github.com/ssjoy/website-health-checker/internal/handlers"
	"github.com/ssjoy/website-health-checker/internal/redis"
	"github.com/ssjoy/website-health-checker/internal/routes"
	"github.com/ssjoy/website-health-checker/internal/services"

	_ "github.com/ssjoy/website-health-checker/docs"
)

func main() {
	// 1. Initialize structured JSON logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	logger.Info("starting health checker application")

	// 2. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load configuration", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Info("configuration loaded successfully", 
		slog.String("app_name", cfg.AppName),
		slog.String("env", cfg.AppEnv),
		slog.Int("port", cfg.AppPort),
	)

	// 3. Connect to MySQL Database (with pooling)
	logger.Info("connecting to mysql database", slog.String("host", cfg.DBHost), slog.Int("port", cfg.DBPort))
	db, err := database.Connect(cfg)
	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		logger.Info("closing mysql database connection pool")
		if err := db.Close(); err != nil {
			logger.Error("error closing mysql database pool", slog.Any("error", err))
		}
	}()
	logger.Info("mysql connection pool initialized successfully")

	// 4. Connect to Redis Cache (with pooling)
	logger.Info("connecting to redis cache", slog.String("host", cfg.RedisHost), slog.Int("port", cfg.RedisPort))
	rdb, err := redis.Connect(cfg)
	if err != nil {
		logger.Error("failed to connect to redis cache", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		logger.Info("closing redis cache connection client")
		if err := rdb.Close(); err != nil {
			logger.Error("error closing redis cache client", slog.Any("error", err))
		}
	}()
	logger.Info("redis client pool initialized successfully")

	// 5. Initialize Services, Handlers and Router
	startTime := time.Now()
	healthService := services.NewHealthService(db, rdb, startTime)
	cacheService := services.NewCacheService(rdb)

	healthHandler := handlers.NewHealthHandler(healthService)
	cacheHandler := handlers.NewCacheHandler(cacheService)

	router := routes.Setup(logger, healthHandler, cacheHandler)

	// 6. Configure HTTP Server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.AppPort),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 7. Start HTTP Server in background goroutine
	go func() {
		logger.Info("starting http server", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server failed to serve", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// 8. Graceful Shutdown signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	sig := <-stop
	logger.Info("received shutdown signal, starting graceful shutdown", slog.String("signal", sig.String()))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", slog.Any("error", err))
	} else {
		logger.Info("http server shut down cleanly")
	}
}
