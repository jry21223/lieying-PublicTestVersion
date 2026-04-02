package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lieying/engine/internal/app"
	"github.com/lieying/engine/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}

	go func() {
		if err := application.Start(); err != nil {
			log.Fatalf("Failed to start app: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	if err := application.Stop(); err != nil {
		log.Fatalf("Failed to stop app: %v", err)
	}
}
