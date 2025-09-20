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
	"github.com/shanto-323/Chat-Server-1/message-service/internal/broker"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/database"
	"github.com/tinrab/retry"
)

type config struct {
	ScyllaUrl string `envconfig:"SCYLLA_URL"`
	RabbitUrl string `envconfig:"RABBIT_URL"`
}

func main() {
	var cfg config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Panic(err)
	}

	var (
		err  error
		conn *amqp.Connection
		repo database.MessageRepository
	)
	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			conn, err = broker.RabbitConnection(cfg.RabbitUrl)
			if err != nil {
				slog.Error(err.Error())
				return err
			}

			repo, err = database.NewUserRepository(cfg.ScyllaUrl)
			if err != nil {
				slog.Error(err.Error())
				return err
			}
			return nil
		},
	)

	br, err := broker.NewMessageBroker(conn)
	if err != nil {
		log.Panic(err)
	}

	slog.Info("MESSAGE-SERVICE RUNNING")

	brokerCtx, brokerCancel := context.WithCancel(context.Background())
	consumer := broker.NewConsumer(br)

	errChan := make(chan error, 1)
	delivery, err := consumer.Consume()
	if err != nil {
		errChan <- err
		slog.Error("MAIN", "consumer", err.Error())
	}
	service := database.NewUserService(repo)

	publisher := broker.NewPublisher(br, service)
	go publisher.Publish(brokerCtx)

	go func() {
		for {
			select {
			case <-brokerCtx.Done():
				return
			case d := <-delivery:
				slog.Info(d.AppId)
				publisher.Delevery <- d
			}
		}
	}()

	stopChen := make(chan os.Signal, 1)
	signal.Notify(stopChen, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stopChen:
		brokerCancel()
		slog.Info("closing server...")
	case err := <-errChan:
		slog.Error(err.Error())
	}
}
