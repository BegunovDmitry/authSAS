package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"authSAS/internal/app"
	"authSAS/internal/config"
	"authSAS/internal/services"
	"authSAS/internal/storages/mockups"
	"authSAS/internal/storages/postgres"
	redisStorage "authSAS/internal/storages/redis"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	testMode = "test"
	localMode = "local"
	productionMode = "prod"
)

func main() {
	startServer := time.Now()
	cfg := config.MustLoad()

	logger := initLogger(cfg, os.Stdout)

	ctx := context.Background()

	logger.Info("Server starting...")
	startApp := time.Now()

	pool, permanentStorage := initPermanentStorage(logger, cfg, ctx)
	if pool != nil {
		defer pool.Close()
	}

	client, temporaryStorage := initTemporaryStorage(logger, cfg)
	if client != nil {
		defer client.Close()
	}
	
	application := app.NewApp(logger, cfg, permanentStorage, temporaryStorage)

	logger.Info("Application initialized", "op_time", time.Since(startApp).Milliseconds())

	go func() {
		application.MustRun()
	}()

	logger.Info(fmt.Sprintf("Server started on port - :%d", cfg.Grpc.Port), "op_time", time.Since(startServer).Milliseconds())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	stopSignal := <- stop

	logger.Info("Server stopping...", "signal", stopSignal.String())
	stopServer := time.Now()

	application.StopApp()

	logger.Info("Server stopped", "op_time", time.Since(stopServer).Milliseconds())
}

func initLogger(cfg *config.Config, output io.Writer) *slog.Logger {
	var handler slog.Handler

	switch cfg.AppMode {
	case testMode:
		handler = slog.NewTextHandler(output, &slog.HandlerOptions{Level: slog.LevelDebug})
	case localMode:
		handler = slog.NewTextHandler(output, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug})
	case productionMode:
		handler = slog.NewJSONHandler(output, &slog.HandlerOptions{Level: slog.LevelDebug})
	}

	logger := slog.New(handler).With(slog.String("app_mode", cfg.AppMode))

	return logger
}

func initPermanentStorage(logger *slog.Logger, cfg *config.Config, ctx context.Context) (*pgxpool.Pool, services.PermanentStorage){
	logger.Info("Permanent DB initialization")
	start := time.Now()

	var pool *pgxpool.Pool
	var permanentStorage services.PermanentStorage

	switch cfg.AppMode {
	case testMode:
		permanentStorage = mockups.NewPermStorMokup()

	case localMode:
		pool, err := pgxpool.New(ctx, cfg.PermStoragePath)
		if err != nil {
			panic(`permanent db pool init error:`)
		}
		permanentStorage = postgres.NewStorage(pool)

	case productionMode:
		pool, err := pgxpool.New(ctx, cfg.PermStoragePath)
		if err != nil {
			panic(`permanent db pool init error:`)
		}
		permanentStorage = postgres.NewStorage(pool)
	}

	logger.Info("Permanent DB initialized", "op_time", time.Since(start).Milliseconds())

	return pool, permanentStorage
}

func initTemporaryStorage(logger *slog.Logger, cfg *config.Config) (*redis.Client, services.TemporaryStorage) {
	logger.Info("Temporary DB initialization")
	start := time.Now()

	var client *redis.Client
	var temporaryStorage services.TemporaryStorage

	switch cfg.AppMode {
	case testMode:
		temporaryStorage = mockups.NewTempStorMokup()

	case localMode:
		opt, err := redis.ParseURL(cfg.TempStorage.TempStoragePath)
		if err != nil {
			panic(`redis init error:`)
		}
		client := redis.NewClient(opt)
		temporaryStorage = redisStorage.NewStorage(client, cfg.TempStorage.CodeTTL)

	case productionMode:
		opt, err := redis.ParseURL(cfg.TempStorage.TempStoragePath)
		if err != nil {
			panic(`redis init error:`)
		}
		client := redis.NewClient(opt)
		temporaryStorage = redisStorage.NewStorage(client, cfg.TempStorage.CodeTTL)
	}


	logger.Info("Temporary DB initialized", "op_time", time.Since(start).Milliseconds())

	return client, temporaryStorage
}