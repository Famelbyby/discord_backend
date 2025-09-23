package app

import (
	mongo_db "discord_backend/configs/mongo"
	common "discord_backend/internal/app"
	grpcApp "discord_backend/internal/app/grpc/relations"
	"discord_backend/internal/services/relations"
	repositoryRelations "discord_backend/internal/storage/relations"
	"fmt"
	"log/slog"
)

type App struct {
	GRPCServer *common.App
}

func New(
	log *slog.Logger,
	grpcPort int,
) (*App, error) {

	mongoDb, err := mongo_db.Connect()
	if err != nil {
		return nil, fmt.Errorf("[ ERROR ] не инициализируется relations  %v", err)
	}

	repo := repositoryRelations.NewStorage(mongoDb, "relations", log)

	profilesService := relations.New(log, repo)
	grpcApp := grpcApp.New(log, profilesService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}, nil
}
