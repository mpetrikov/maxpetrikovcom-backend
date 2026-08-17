package worker

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"
)

func (w *Worker) Run() error {
	defer func() {
		w.rabbitMQ.Close()
		w.db.Close()

		w.logger.Info("worker resources closed")
	}()

	deliveries, err := w.rabbitMQ.ConsumeLabCreate()
	if err != nil {
		return fmt.Errorf(
			"start lab.create consumer: %w",
			err,
		)
	}

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	w.logger.Info("worker started")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("worker shutdown signal received")
			return nil

		case delivery, ok := <-deliveries:
			if !ok {
				return errors.New(
					"lab.create delivery channel closed",
				)
			}

			if err := w.handleLabCreate(
				ctx,
				delivery.Body,
			); err != nil {
				w.logger.Error(
					"failed to process lab.create",
					"error", err,
				)

				_ = delivery.Nack(
					false,
					true,
				)

				continue
			}

			_ = delivery.Ack(false)
		}
	}
}
