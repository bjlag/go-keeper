package main

import (
	"context"
	logNative "log"
	"os/signal"
	"syscall"

	"github.com/bjlag/go-keeper/internal/app"
	"github.com/bjlag/go-keeper/internal/infrastructure/logger"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logNative.Fatalf("panic occurred: %v", r)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	log := logger.Get("dev")
	defer func() {
		_ = log.Sync()
	}()

	err := app.NewApp(log).Run(ctx)
	if err != nil {
		panic(err)
	}
}
