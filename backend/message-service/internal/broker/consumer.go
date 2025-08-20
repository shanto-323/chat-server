package broker

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Consume() (<-chan amqp.Delivery, error)
}

type consumer struct {
	messageBroker MessageBroker
	ctx           context.Context
}

func NewConsumer(ctx context.Context, broker MessageBroker) Consumer {
	return &consumer{
		messageBroker: broker,
		ctx:           ctx,
	}
}

func (c *consumer) Consume() (<-chan amqp.Delivery, error) {
	err := c.messageBroker.CreateQueue(MESSAGE_QUEUE, true, false)
	if err != nil {
		return nil, err
	}
	err = c.messageBroker.CreateQueueBinding(MESSAGE_QUEUE, ROUTING_KEY_MESSAGE, EXCHANGE_KEY)
	if err != nil {
		return nil, err
	}
	return c.messageBroker.Consume(MESSAGE_QUEUE, "", false)
}
