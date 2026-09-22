package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/gbadopt/gbadopt/internal/config"
	"github.com/gbadopt/gbadopt/internal/router"
	"github.com/gbadopt/gbadopt/internal/util"
)

func main() {
	logger := util.NewLogger()
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		os.Exit(1)
	}
	if err := migrate(db); err != nil {
		logger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}
	if err := seed(db); err != nil {
		logger.Error("failed to seed database", "error", err)
		os.Exit(1)
	}

	redis, err := util.NewRedisClient(cfg)
	if err != nil {
		logger.Error("failed to connect redis", "error", err)
		os.Exit(1)
	}
	minio, err := util.NewMinIOClient(cfg)
	if err != nil {
		logger.Error("failed to connect minio", "error", err)
		os.Exit(1)
	}

	r := router.Setup(cfg, db, redis, minio, logger)

	srv := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("gbadopt backend listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server stopped", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", "error", err)
	}
}
