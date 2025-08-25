package broker

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	EXCHANGE_KEY  = "message.service"
	EXCHANGE_TYPE = "topic"

	MESSAGE_QUEUE = "message.queue"

	ROUTING_KEY_MESSAGE = "incomming.message"
)

type MessageBroker interface {
	Close() error
	CreateQueue(queueName string, durable, autoDelete bool) error
	CreateQueueBinding(name, binding, exchange string) error
	SendMessage(ctx context.Context, exchange, routingKey string, opt amqp.Publishing) error
	Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error)
}

type rabbitClient struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func RabbitConnection(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}

func NewMessageBroker(conn *amqp.Connection) (MessageBroker, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		"message.service",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		return nil, err
	}

	return &rabbitClient{
		conn: conn,
		ch:   ch,
	}, nil
}

func (rc *rabbitClient) Close() error {
	if rc.ch != nil {
		rc.ch.Close()
	}
	return rc.conn.Close()
}

func (rc *rabbitClient) CreateQueue(queueName string, durable, autoDelete bool) error {
	_, err := rc.ch.QueueDeclare(queueName, durable, autoDelete, false, false, nil)
	return err
}

func (rc *rabbitClient) CreateQueueBinding(name, binding, exchange string) error {
	return rc.ch.QueueBind(name, binding, exchange, false, nil)
}

func (rc *rabbitClient) SendMessage(ctx context.Context, exchange, routingKey string, opt amqp.Publishing) error {
	return rc.ch.PublishWithContext(ctx, exchange, routingKey, true, false, opt)
}

func (rc *rabbitClient) Consume(queue, consumer string, autoAck bool) (<-chan amqp.Delivery, error) {
	return rc.ch.Consume(queue, consumer, autoAck, false, false, false, nil)
}
