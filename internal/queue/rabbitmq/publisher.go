package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/job"
)

func (c *Client) PublishCreate(
	ctx context.Context,
	job job.LabCreate,
) error {
	return c.publishJSON(
		ctx,
		LabCreateQueue,
		job,
	)
}

func (c *Client) PublishStop(
	ctx context.Context,
	job job.LabStop,
) error {
	return c.publishJSON(
		ctx,
		LabStopQueue,
		job,
	)
}

func (c *Client) publishJSON(
	ctx context.Context,
	queue string,
	message any,
) error {
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal %s job: %w", queue, err)
	}

	c.mu.RLock()
	channel := c.channel
	c.mu.RUnlock()

	if channel == nil || channel.IsClosed() {
		return errors.New("RabbitMQ channel is unavailable")
	}

	err = channel.PublishWithContext(
		ctx,
		"",    // default exchange
		queue, // routing key == queue name
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish %s job: %w", queue, err)
	}

	return nil
}
