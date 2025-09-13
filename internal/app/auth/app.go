package app

import (
	common "idea-store-auth/internal/app"
	grpcApp "idea-store-auth/internal/app/grpc/auth"
	"idea-store-auth/internal/services/auth"
	"idea-store-auth/internal/storage/postgre"
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
