package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/api"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/connection"
)

type config struct {
	GatewayPort string `envconfig:"GATEWAY_PORT"`
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	manager := connection.NewManager(ctx, nil)
	api := api.NewApi(cfg.GatewayPort, manager)

	stopChan := make(chan os.Signal, 1)
	go func() {
		if err := api.Start(); err != nil {
			slog.Error(err.Error())
			stopChan <- syscall.SIGINT
		}
	}()

	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
	slog.Info("Closing the server....")
}
