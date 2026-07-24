package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/artcodefun/detective-game/backend/internal/bootstrap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := bootstrap.LoadConfig()
	module := bootstrap.NewModule(cfg)

	if err := module.Run(ctx); err != nil {
		log.Fatalf("Server exited with error: %v", err)
	}
}
