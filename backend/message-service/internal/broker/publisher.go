package broker

import (
	"context"
	"encoding/json"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/broker/model"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/database"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/remote"
)

type Publisher interface {
	Publish(d amqp.Delivery) error
}

type publisher struct {
	messageBroker MessageBroker
	cacheClient   remote.CacheClient
	service       *database.MessageService
}

func NewPublisher(broker MessageBroker, service *database.MessageService) Publisher {
	return &publisher{
		messageBroker: broker,
		cacheClient:   remote.NewCacheClient(),
		service:       service,
	}
}

func (p *publisher) Publish(d amqp.Delivery) error {
	packet := model.Packet{}
	if err := json.Unmarshal(d.Body, &packet); err != nil {
		return err
	}

	switch packet.Type {
	case model.TYPE_CHAT:
		p.handleTypeChat(&packet)
	case model.TYPE_LIST:
		return p.handleTypeList(packet.ReceiverId)
	}

	slog.Info("BROKER", "success", packet)
	return nil
}

func (p *publisher) brodcast(uid string, payload *model.GatewayPayload) error {
	resp, err := p.cacheClient.GetActivePool(uid)
	if err != nil {
		return err
	}

	for id, value := range resp.Message.ActivePool {
		payload.SessionId = id

		body, err := json.Marshal(payload)
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

	return nil
}

func (p *publisher) handleTypeList(uid string) error {
	resp, err := p.cacheClient.GetAlConnPool()
	if err != nil {
		return err
	}

	payload := model.GatewayPayload{
		Pool: resp.Message.ConnPool,
	}
	return p.brodcast(uid, &payload)
}

func (p *publisher) handleTypeChat(packet *model.Packet) error {
	// CHECK IF USER ONLINE
	resp, err := p.cacheClient.GetActivePool(packet.ReceiverId)
	if err != nil {
		return err
	}

	// USER IS ONLINE // marge in a switch statement
	offline := true
	if resp.Status == 200 {
		offline = false
		if err := p.brodcast(packet.ReceiverId, &model.GatewayPayload{Data: packet.Data}); err != nil {
			return err
		}
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := p.service.PushMessage(ctx, packet.SenderId, packet.ReceiverId, packet.Data, offline); err != nil {
		return err
	}
	return nil
}
