package app

import (
	common "discord_backend/internal/app"
	grpcApp "discord_backend/internal/app/grpc/auth"
	"discord_backend/internal/services/auth"
	"discord_backend/internal/storage/postgre"
	"discord_backend/internal/storage/redis"
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
	postgresStorage, err := postgre.New()

	if err != nil {
		panic(err)
	}

	redisStorage, err := redis.NewClient()

	if err != nil {
		panic(err)
	}

	authService := auth.New(log, postgresStorage, postgresStorage, redisStorage, tokenTTL)

	grpcApp := grpcApp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
