package rabbitmq

import (
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const LabCreateQueue = "lab.create"
const maxReconnectDelay = 5 * time.Second

type Client struct {
	url string

	mu         sync.RWMutex
	connection *amqp.Connection
	channel    *amqp.Channel

	closeCh chan struct{}
}

func New(
	url string,
) (*Client, error) {
	client := &Client{
		url:     url,
		closeCh: make(chan struct{}),
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	go client.watchConnection()

	return client, nil
}

func (c *Client) connect() error {
	connection, err := amqp.Dial(c.url)
	if err != nil {
		return fmt.Errorf("connect to RabbitMQ: %w", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		_ = connection.Close()
		return fmt.Errorf("open RabbitMQ channel: %w", err)
	}

	if err := c.declareQueues(channel); err != nil {
		_ = channel.Close()
		_ = connection.Close()

		return err
	}

	c.mu.Lock()

	oldChannel := c.channel
	oldConnection := c.connection

	c.connection = connection
	c.channel = channel

	c.mu.Unlock()

	if oldChannel != nil {
		_ = oldChannel.Close()
	}

	if oldConnection != nil {
		_ = oldConnection.Close()
	}

	return nil
}

func (c *Client) watchConnection() {
	for {
		c.mu.RLock()
		connection := c.connection
		c.mu.RUnlock()

		if connection == nil {
			c.reconnect()
			continue
		}

		notifyClose := connection.NotifyClose(
			make(chan *amqp.Error, 1),
		)

		select {
		case <-c.closeCh:
			return

		case <-notifyClose:
			c.reconnect()
		}
	}
}

func (c *Client) reconnect() {
	delay := time.Second

	for {
		select {
		case <-c.closeCh:
			return

		case <-time.After(delay):
		}

		if err := c.connect(); err == nil {
			return
		}

		delay *= 2
		if delay > maxReconnectDelay {
			delay = maxReconnectDelay
		}
	}
}

func (c *Client) declareQueues(channel *amqp.Channel) error {
	_, err := channel.QueueDeclare(
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
	select {
	case <-c.closeCh:
		return
	default:
		close(c.closeCh)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.channel != nil {
		_ = c.channel.Close()
		c.channel = nil
	}

	if c.connection != nil {
		_ = c.connection.Close()
		c.connection = nil
	}
}
