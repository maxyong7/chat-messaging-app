package main

import (
	"log"

	"github.com/maxyong7/chat-messaging-app/config"
	"github.com/maxyong7/chat-messaging-app/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
