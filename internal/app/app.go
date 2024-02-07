package app

import (
	grpcapp "github.com/ARUMANDESU/uniclubs-filestorage-service/internal/app/grpc"
	httpapp "github.com/ARUMANDESU/uniclubs-filestorage-service/internal/app/http"
	"github.com/ARUMANDESU/uniclubs-filestorage-service/internal/config"
	"github.com/ARUMANDESU/uniclubs-filestorage-service/internal/handler"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
	HTTPSrv *httpapp.App
}

func New(log *slog.Logger, cfg *config.Config) *App {
	const op = "App.New"
	_ = log.With(slog.String("op", op))

	grpcApp := grpcapp.New(log, cfg.GRPC.Port)

	h := handler.New(log)
	httpApp := httpapp.New(cfg, log, h.InitRoutes())

	return &App{GRPCSrv: grpcApp, HTTPSrv: httpApp}
}
