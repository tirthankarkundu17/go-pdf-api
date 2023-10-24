package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/tirthankarkundu17/pdf-api/api"
	"go.uber.org/zap"
)

var (
	logger, _ = zap.NewProduction(zap.Fields(zap.String("type", "main")))
	shutdowns []func() error
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	var (
		ctx     = context.Background()
		handler = api.New()
		server  = http.Server{
			Addr:    fmt.Sprintf("%s:%s", host, port),
			Handler: handler,
		}
		shutdown = make(chan struct{})
	)

	go gracefulShutdown(ctx, &server, shutdown)

	logger.Info("server starting: " + server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("server error", zap.Error(err))
	}

	<-shutdown
}

func gracefulShutdown(ctx context.Context, server *http.Server, shutdown chan struct{}) {
	var (
		sigint = make(chan os.Signal, 1)
	)

	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	logger.Info("shutting down server gracefully")

	// stop receiving any request.
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("shutdown error", zap.Error(err))
	}

	// close any other modules.
	for i := range shutdowns {
		shutdowns[i]()
	}

	close(shutdown)
}
