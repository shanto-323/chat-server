package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shanto-323/Chat-Server-1/message-service/cmd/model"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/queue"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/remote"
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
			packet := model.Packet{}
			if err := json.Unmarshal(m.Body, &packet); err != nil {
				slog.Error(err.Error())
				continue
			}

			cacheClient := remote.NewCacheClient()
			resp, err := cacheClient.GetActivePool(packet.ReceiverId)
			if err != nil {
				slog.Error(err.Error())
				continue
			}

			for id, value := range resp.Message.ActivePool {
				publishPacket := model.PublishPacket{
					SessionId: id,
					Data:      packet.Data,
				}

				body, err := json.Marshal(&publishPacket)
				if err != nil {
					slog.Error(err.Error())
					continue
				}

				if err := queue.SendMessage(context.Background(), "message.service", value, amqp091.Publishing{
					ContentType:  "application/json",
					DeliveryMode: amqp091.Persistent,
					Body:         body,
				}); err != nil {
					slog.Error(err.Error())
					continue
				}
				log.Println(string(m.Body))
			}
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
