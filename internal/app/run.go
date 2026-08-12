package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
)

func (a *App) Run() error {
	serverError := make(chan error, 1)

	go func() {
		a.logger.Info(
			"starting HTTP server",
			"address", a.config.HTTPAddress,
			"environment", a.config.Environment,
		)

		serverError <- a.httpServer.ListenAndServe()
	}()

	shutdownContext, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverError:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf(
				"HTTP server failed: %w",
				err,
			)
		}

	case <-shutdownContext.Done():
		a.logger.Info("shutdown signal received")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		a.config.ShutdownTimeout,
	)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf(
			"shutdown HTTP server: %w",
			err,
		)
	}

	a.db.Close()

	a.logger.Info("HTTP server stopped")

	return nil
}
