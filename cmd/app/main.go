package main

import (
	"github.com/AntonBraer/urlShorter/config"
	"github.com/AntonBraer/urlShorter/internal/app"
	"log"
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
