package main

import (
	"context"
	"errors"
	"github.com/ARUMANDESU/uniclubs-filestorage-service/internal/app"
	"github.com/ARUMANDESU/uniclubs-filestorage-service/internal/config"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("starting file storage service",
		slog.String("env", cfg.Env),
		slog.String("address", cfg.HTTP.Address),
	)

	application := app.New(log, cfg)
	go application.GRPCSrv.MustRun()

	go func() {
		if err := application.HTTPSrv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server error: %v", slog.Attr{
				Key:   "error",
				Value: slog.StringValue(err.Error()),
			})
		}
		log.Info("stopped serving new http connections")

	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := application.HTTPSrv.Stop(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("shutdown error", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}

	log.Info("application stopped")

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
