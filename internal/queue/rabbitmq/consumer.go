package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Client) ConsumeLabCreate() (
	<-chan amqp.Delivery,
	error,
) {
	return c.consumeQueue(LabCreateQueue)
}

func (c *Client) ConsumeLabStop() (
	<-chan amqp.Delivery,
	error,
) {
	return c.consumeQueue(LabStopQueue)
}

func (c *Client) consumeQueue(
	queue string,
) (<-chan amqp.Delivery, error) {
	c.mu.RLock()
	channel := c.channel
	c.mu.RUnlock()

	if channel == nil || channel.IsClosed() {
		return nil, fmt.Errorf(
			"RabbitMQ channel is unavailable",
		)
	}

	deliveries, err := channel.Consume(
		queue,
		"",
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"consume %s queue: %w",
			queue,
			err,
		)
	}

	return deliveries, nil
}
