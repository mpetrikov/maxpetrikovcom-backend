package app

import (
	"fmt"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/queue/rabbitmq"
)

func buildRabbitMQ(
	url string,
) (*rabbitmq.Client, error) {
	client, err := rabbitmq.New(url)
	if err != nil {
		return nil, fmt.Errorf("initialize RabbitMQ: %w", err)
	}

	return client, nil
}
