package broker

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/broker/model"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/database"
	"github.com/shanto-323/Chat-Server-1/message-service/internal/remote"
)

type Publisher struct {
	messageBroker MessageBroker
	Delevery      chan amqp.Delivery
	cacheClient   remote.CacheClient
	service       *database.MessageService
	activePool    map[string]bool
}

func NewPublisher(broker MessageBroker, service *database.MessageService) *Publisher {
	return &Publisher{
		messageBroker: broker,
		Delevery:      make(chan amqp.Delivery),
		cacheClient:   remote.NewCacheClient(),
		service:       service,
		activePool:    map[string]bool{},
	}
}

func (p *Publisher) Publish(ctx context.Context) {
	go p.resetPool(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case d := <-p.Delevery:
			go func() {
				if err := p.process(d); err != nil {
					slog.Error("BROKER", "Publish", err.Error())
				}
			}()
		}
	}
}

func (p *Publisher) resetPool(ctx context.Context) {
	hold := time.NewTicker(2 * time.Second)
	defer hold.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-hold.C:
			resp, err := p.cacheClient.GetAlConnPool()
			if err != nil {
				slog.Error("BROKER", "resetPool", err.Error())
				time.Sleep(2 * time.Second)
				continue
			}

			tempPool := make(map[string]bool)
			for _, conn := range resp.Message.ConnPool {
				tempPool[conn] = true
			}

			p.activePool = tempPool
			time.Sleep(2 * time.Second)
		}
	}
}

func (p *Publisher) process(d amqp.Delivery) error {
	packet := model.Packet{}
	if err := json.Unmarshal(d.Body, &packet); err != nil {
		return err
	}

	switch packet.Type {
	case model.TYPE_CHAT:
		{
			incommingMessage := model.IncommingMessage{}
			if err := json.Unmarshal(packet.Payload, &incommingMessage); err != nil {
				return err
			}
			slog.Info("CHAT", "incommingmsg", incommingMessage)

			// CHECK IF USER ONLINE ONLINE
			_, offline := p.activePool[packet.ReceiverId]

			// STORE IN DATABASE
			if err := p.service.PushMessage(
				context.Background(), // NEED TO WORK (maybe wg)
				packet.SenderId,
				packet.ReceiverId,
				incommingMessage.Message, // FROM RAW MESSAGE
				offline,
			); err != nil {
				return err
			}

			// REALTIME LAST 10 MESSAGE FROM BOTH END
			messages, err := p.service.GetLatestMessage(context.Background(), packet.SenderId, packet.ReceiverId)
			if err != nil {
				return err
			}
			clientMsg := model.Messages{}
			for _, m := range messages {
				clientMsg.Messages = append(clientMsg.Messages, *m)
			}
			rawMsg, err := json.Marshal(&clientMsg)
			if err != nil {
				slog.Error("BROKER", "json", err.Error())
				return err
			}

			eventPacket := model.EventPacket{
				Type:    model.TYPE_CHAT,
				Payload: rawMsg,
			}

			go func() {
				p.brodcast(packet.ReceiverId, &eventPacket)
			}()

			// FOR SENDER UPDATE
			go func() {
				p.brodcast(packet.SenderId, &eventPacket)
			}()

		}
	}

	return nil
}

func (p *Publisher) brodcast(id string, eventPacket *model.EventPacket) {
	_, exists := p.activePool[id]
	if !exists {
		slog.Info("BROKER", "user not active", id)
		return
	}

	resp, err := p.cacheClient.GetActivePool(id)
	if err != nil {
		slog.Error("BROKER", "brodcast", err.Error())
	}

	for sessionId, gatewayId := range resp.Message.ActivePool {
		pack := eventPacket
		pack.SessionId = sessionId

		body, err := json.Marshal(&pack)
		if err != nil {
			slog.Error("BROKER", "brodcast", err.Error())
			continue
		}

		slog.Info("RESP", "info pack", string(eventPacket.Payload))
		p.messageBroker.SendMessage(context.Background(), EXCHANGE_KEY, gatewayId, amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		})
	}
}
