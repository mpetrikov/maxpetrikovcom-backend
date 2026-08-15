package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

const LabCreateQueue = "lab.create"

type Client struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func New(
	url string,
) (*Client, error) {
	connection, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		connection.Close()
		return nil, fmt.Errorf("open RabbitMQ channel: %w", err)
	}

	client := &Client{
		connection: connection,
		channel:    channel,
	}

	if err := client.declareQueues(); err != nil {
		client.Close()
		return nil, err
	}

	return client, nil
}

func (c *Client) declareQueues() error {
	_, err := c.channel.QueueDeclare(
		LabCreateQueue,
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("declare %s queue: %w", LabCreateQueue, err)
	}

	return nil
}

func (c *Client) Close() {
	if c.channel != nil {
		_ = c.channel.Close()
	}

	if c.connection != nil {
		_ = c.connection.Close()
	}
}
