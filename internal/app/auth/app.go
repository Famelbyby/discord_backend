package app

import (
	common "discord_backend/internal/app"
	grpcApp "discord_backend/internal/app/grpc/auth"
	"discord_backend/internal/services/auth"
	"discord_backend/internal/storage/postgre"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *common.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenTTL time.Duration,
) *App {
	storage, err := postgre.New()
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, tokenTTL)

	grpcApp := grpcApp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
