package rabbitmq

import (
	"fmt"
	"log/slog"

	"github.com/streadway/amqp"
)

type RabbitClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

func NewRabbitClient(connectionUrl string) (*RabbitClient, error) {
	conn, channel, err := newConnection(connectionUrl)
	if err != nil {
		return nil, err
	}
	return &RabbitClient{
		conn:    conn,
		channel: channel,
		url:     connectionUrl,
	}, nil
}

func (r *RabbitClient) ConsumeMessages(exchangeName, routingKey, queueName string) (<-chan amqp.Delivery, error) {
	err := r.channel.ExchangeDeclare(exchangeName, "direct", true, true, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}
	queue, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}
	err = r.channel.QueueBind(queue.Name, routingKey, exchangeName, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %v", err)
	}
	msgs, err := r.channel.Consume(queue.Name, "goapp", false, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to consume messages from queue: %v", err)
	}
	slog.Info("Connected to RabbitMQ successfully")
	return msgs, nil
}

func (r *RabbitClient) PublishMessage(exchangeName, routingKey, queueName string, message []byte) error {
	err := r.channel.ExchangeDeclare(exchangeName, "direct", true, true, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %v", err)
	}
	queue, err := r.channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}
	err = r.channel.QueueBind(queue.Name, routingKey, exchangeName, false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}
	err = r.channel.Publish(exchangeName, routingKey, false, false, amqp.Publishing{ContentType: "application/json", Body: message})
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}
	return nil
}

func (r *RabbitClient) Close() {
	r.channel.Close()
	r.conn.Close()
}

func newConnection(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open a channel: %v", err)
	}
	return conn, channel, nil
}
