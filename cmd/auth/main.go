package main

import (
	"os"
	"os/signal"
	"syscall"

	common "idea-store-auth/cmd"
	appAuth "idea-store-auth/internal/app/auth"
	"idea-store-auth/internal/config"
)

func main() {
	cfg := config.MustLoad()

	log := common.SetupLogger(cfg.Env)

	authApp := appAuth.New(log, cfg.GRPC.AuthMS.Port, cfg.TokenTTL)

	go func() {
		authApp.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	authApp.GRPCServer.Stop()
	log.Info("Gracefully stopped")
}
