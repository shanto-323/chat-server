package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/queue"
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

	queue, err := queue.NewConsumer(conn)
	if err != nil {
		log.Panic(err)
	}

	errChan := make(chan error, 1)
	slog.Info("MESSAGE-SERVICE RUNNING")
	go func() {
		err := queue.CreateQueue("message.queue", true, false)
		if err != nil {
			slog.Error(err.Error())
			return
		}
		err = queue.CreateQueueBinding("message.queue", "incomming.message", "message.service")
		if err != nil {
			slog.Error(err.Error())
			return
		}
		msg, err := queue.Consume("message.queue", "", false)
		if err != nil {
			slog.Error(err.Error())
			return
		}

		for m := range msg {
			log.Println(string(m.Body))
		}
	}()

	stopChen := make(chan os.Signal, 1)
	signal.Notify(stopChen, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stopChen:
		slog.Info("closing server...")
	case err := <-errChan:
		slog.Error(err.Error())
	}
}
