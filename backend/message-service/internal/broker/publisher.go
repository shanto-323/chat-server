package broker

import (
	"context"
	"encoding/json"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/broker/model"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/remote"
)

type Publisher interface {
	Publish(d amqp.Delivery) error
}

type publisher struct {
	messageBroker MessageBroker
}

func NewPublisher(broker MessageBroker) Publisher {
	return &publisher{
		messageBroker: broker,
	}
}

func (p *publisher) Publish(d amqp.Delivery) error {
	packet := model.Packet{}
	if err := json.Unmarshal(d.Body, &packet); err != nil {
		return err
	}

	cacheClient := remote.NewCacheClient()
	resp, err := cacheClient.GetActivePool(packet.ReceiverId)
	if err != nil {
		return err
	}

	for id, value := range resp.Message.ActivePool {
		publishPacket := model.GatewayPayload{
			SessionId: id,
			Data:      packet.Data,
		}

		body, err := json.Marshal(&publishPacket)
		if err != nil {
			return err
		}

		if err := p.messageBroker.SendMessage(context.Background(), EXCHANGE_KEY, value, amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		}); err != nil {
			return err
		}
	}

	slog.Info("BROKER", "success", packet)
	return nil
}
