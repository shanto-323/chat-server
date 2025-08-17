package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/api"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/connection"
	"github.com/shanto-323/Chat-Server-1/gateway-1/pkg/queue"
	"github.com/tinrab/retry"
)

type config struct {
	GatewayPort string `envconfig:"GATEWAY_PORT"`
	RabbitUrl   string `envconfig:"RABBIT_URL"`
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Panic(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		err  error
		conn *amqp.Connection
	)
	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			conn, err = queue.RabbitConnection(cfg.RabbitUrl)
			if err != nil {
				slog.Error(err.Error())
				return err
			}
			return nil
		},
	)

	consumer, err := queue.NewConsumer(conn)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	manager := connection.NewManager(ctx, consumer)
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
