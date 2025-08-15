package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/shanto-323/Chat-Server-1/client-service/internal/api"
	"github.com/shanto-323/Chat-Server-1/client-service/internal/database"
	"github.com/tinrab/retry"
)

type config struct {
	ClientServicePort string `envconfig:"CLIENT_SERVICE_PORT"`
	PostgresUrl       string `envconfig:"POSTGRES_URL"`
	ReddisUrl         string `envconfig:"REDIS_URL"`
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Panic(err)
	}

	var (
		err        error
		repository database.UserRepository
	)
	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			repository, err = database.NewUserRepository(cfg.PostgresUrl)
			if err != nil {
				slog.Error(err.Error())
				return err
			}
			return nil
		},
	)

	service := database.NewUserService(repository)
	api := api.NewApi(cfg.ClientServicePort, nil, service)

	errChan := make(chan error, 1)
	slog.Info("SERVER RUNNING", "client-service PORT:", cfg.ClientServicePort)
	go func() {
		if err := api.Start(); err != nil {
			errChan <- err
			return
		}
	}()

	stopChen := make(chan os.Signal, 1)
	signal.Notify(stopChen, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case _ = <-stopChen:
		slog.Info("closing server...")
	case err := <-errChan:
		slog.Error(err.Error())
	}
}
