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
	body, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal lab create job: %w", err)
	}

	c.mu.RLock()
	channel := c.channel
	c.mu.RUnlock()

	if channel == nil || channel.IsClosed() {
		return errors.New("RabbitMQ channel is unavailable")
	}

	err = channel.PublishWithContext(
		ctx,
		"",             // default exchange
		LabCreateQueue, // routing key == queue name
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish lab create job: %w", err)
	}

	return nil
}
