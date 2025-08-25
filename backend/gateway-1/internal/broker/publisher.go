package broker

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	Publish(payload []byte) error
}

type publisher struct {
	messageBroker MessageBroker
	Delevery      chan amqp.Delivery
}

func NewPublisher(broker MessageBroker) Publisher {
	return &publisher{
		messageBroker: broker,
		Delevery:      make(chan amqp.Delivery),
	}
}

func (p *publisher) Publish(payload []byte) error {
	return p.messageBroker.SendMessage(context.Background(), EXCHANGE_KEY, ROUTING_KEY_INCOMMING, amqp.Publishing{
		ContentType:  "text/plain",
		DeliveryMode: amqp.Persistent,
		Body:         payload,
	})
}
