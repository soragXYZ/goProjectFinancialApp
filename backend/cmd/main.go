package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"financialApp/api/router"
	"financialApp/config"
)

func main() {

	config.Init()

	router := router.New()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Conf.Server.Port),
		Handler:      router,
		ReadTimeout:  config.Conf.Server.TimeoutRead,
		WriteTimeout: config.Conf.Server.TimeoutWrite,
		IdleTimeout:  config.Conf.Server.TimeoutIdle,
	}

	// Correct way to handle a server shutdown
	// https://dev.to/mokiat/proper-http-shutdown-in-go-3fji

	closed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		config.Logger.Info().Msgf("Shutting down server %v", config.Conf.Server.Port)

		ctx, cancel := context.WithTimeout(context.Background(), server.IdleTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			config.Logger.Error().Err(err).Msg("Server shutdown failure")
		}

		defer config.DB.Close()

		close(closed)
	}()

	config.Logger.Info().Msgf("Starting server on port %v", config.Conf.Server.Port)
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		config.Logger.Fatal().Err(err).Msg("Server failure")
	}

	<-closed
	config.Logger.Info().Msg("Server shutdown successfully")
}
