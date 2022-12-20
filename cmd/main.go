package main

import (
	"auth/internal/server"

	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	server := server.HttpServer{}
	go func() {
		if err := server.Start(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("start server failed: %#v", err)
		}
	}()

	timeoutCtx := 1 * time.Minute
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeoutCtx)
	defer cancel()

	if err := server.Stop(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %#v", err)
	}
}