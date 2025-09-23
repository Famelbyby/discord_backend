package main

import (
	"os"
	"os/signal"
	"syscall"

	common "discord_backend/cmd"
	appRelations "discord_backend/internal/app/relations"
	"discord_backend/internal/config"
)

func main() {
	cfg := config.MustLoad()

	log := common.SetupLogger(cfg.Env)

	relationsApp, err := appRelations.New(log, cfg.GRPC.RelationsMS.Port)

	if err != nil {
		log.Error("[main] error starting relations app: " + err.Error())
		panic("cant start relations app, error: " + err.Error())
	}

	go func() {
		relationsApp.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	relationsApp.GRPCServer.Stop()
	log.Info("Gracefully stopped")
}
