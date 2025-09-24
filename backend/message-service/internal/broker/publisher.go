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
	messageBroker  MessageBroker
	Delevery       chan amqp.Delivery
	cacheClient    remote.CacheClient
	service        *database.MessageService
	activePool     map[string]any
	activePoolList []string
}

func NewPublisher(broker MessageBroker, service *database.MessageService) *Publisher {
	return &Publisher{
		messageBroker:  broker,
		Delevery:       make(chan amqp.Delivery),
		cacheClient:    remote.NewCacheClient(),
		service:        service,
		activePool:     make(map[string]any),
		activePoolList: []string{},
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

			p.activePoolList = resp.Message.ConnPool // INSTENT LIST

			tempPool := make(map[string]any)
			for _, conn := range resp.Message.ConnPool {
				tempPool[conn] = nil
			}

			p.activePool = tempPool // HASH MAP FOR FINDING SINGLE VALUE QUICKLY
			time.Sleep(2 * time.Second)
		}
	}
}

func (p *Publisher) process(d amqp.Delivery) error {
	pWrapper := model.PacketWrapper{}
	if err := json.Unmarshal(d.Body, &pWrapper); err != nil {
		slog.Error("process", "file", string(d.Body))
		return err
	}

	switch pWrapper.Type {
	case model.TYPE_CHAT:
		packet := model.Packet{}
		if err := json.Unmarshal(pWrapper.Payload, &packet); err != nil {
			return err
		}

		_, online := p.activePool[packet.ReceiverId] // CHECK IF USER ONLINE ONLINE
		slog.Info("Process", "Client ->", packet.SenderId, "WANTS TO SEND DATA TO Client ->", packet.ReceiverId, ".ONLINE STATUS ", online)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := p.service.PushMessage( // STORE IN DATABASE
			ctx,
			packet.SenderId,
			packet.ReceiverId,
			packet.Message,
			!online, // IF FOUND offline = !online = !TRUE
		); err != nil {
			return err
		}

		// REALTIME LAST MESSAGE
		message, err := p.service.GetLatestMessage(context.Background(), packet.SenderId, packet.ReceiverId)
		if err != nil {
			return err
		}

		rawMsg, err := json.Marshal(&message)
		if err != nil {
			slog.Error("BROKER", "json", err.Error())
			return err
		}

		eventPacket := model.EventPacket{
			Type:    model.TYPE_CHAT,
			PeerId:  packet.SenderId,
			Payload: rawMsg,
		}

		if online {
			slog.Info("Process", "SEND DATA TO THE RECEIVER ->", packet.ReceiverId, "WITH PEER ID ->", packet.SenderId)
			go p.brodcast(packet.ReceiverId, &eventPacket)
		}

		senderEventPacket := eventPacket
		senderEventPacket.PeerId = packet.ReceiverId
		slog.Info("Process", "SEND DATA TO THE SENDER ->", packet.SenderId, "WITH PEER ID ->", packet.ReceiverId)
		go p.brodcast(packet.SenderId, &senderEventPacket)

	case model.TYPE_CHAT_HISTORY:
		chatHistoryPacket := model.ChatHistoryPacket{}
		if err := json.Unmarshal(pWrapper.Payload, &chatHistoryPacket); err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		messages, err := p.service.GetMessage(ctx, chatHistoryPacket.SenderId, chatHistoryPacket.ReceiverId, chatHistoryPacket.LastUpdate)
		if err != nil {
			return nil
		}

		rawMsg, err := json.Marshal(&messages)
		if err != nil {
			slog.Error("BROKER", "json", err.Error())
			return err
		}

		eventPacket := model.EventPacket{
			Type:    model.TYPE_CHAT,
			Payload: rawMsg,
		}
		go p.brodcast(chatHistoryPacket.SenderId, &eventPacket)
	case model.TYPE_LIST:
		listPacket := model.ListPacket{}
		if err := json.Unmarshal(pWrapper.Payload, &listPacket); err != nil {
			return err
		}
		activePool := model.ActivePool{
			Pool: p.activePoolList,
		}
		rawMsg, err := json.Marshal(&activePool)
		if err != nil {
			slog.Error("BROKER", "json", err.Error())
			return err
		}
		eventPacket := model.EventPacket{
			Type:    model.TYPE_CHAT_HISTORY,
			Payload: rawMsg,
		}
		go p.brodcast(listPacket.Uid, &eventPacket)
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
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("brodcast", "debug peer id", eventPacket.PeerId)

	for sessionId, gatewayId := range resp.Message.ActivePool {
		go func(sId string, ctx context.Context) {
			pack := eventPacket
			pack.SessionId = sId

			body, err := json.Marshal(&pack)
			if err != nil {
				slog.Error("BROKER", "brodcast", err.Error())
				return
			}

			if err := p.messageBroker.SendMessage(ctx, EXCHANGE_KEY, gatewayId, amqp.Publishing{
				ContentType:  "application/json",
				DeliveryMode: amqp.Persistent,
				Body:         body,
			}); err != nil {
				slog.Error("BROKER", "brodcast", err.Error())
				return
			}
		}(sessionId, ctx)

	}
}
