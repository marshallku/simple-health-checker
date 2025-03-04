package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/marshallku/statusy/config"
	"github.com/marshallku/statusy/handler"
	"github.com/marshallku/statusy/health"
	"github.com/marshallku/statusy/store"
)

func main() {
	mode := flag.String("mode", "server", "Mode to run in (server or health)")
	configFile := flag.String("config", "config.yaml", "Path to configuration file")

	flag.Parse()

	var err error
	var cfg *config.Config
	cfg, err = config.LoadConfig(*configFile)

	if err != nil {
		fmt.Printf("Error loading configuration: %v\n", err)
		os.Exit(1)
	}

	if *mode == "cli" {
		health.Check(cfg, nil)
	} else {
		store := store.NewStore()
		server := handler.NewHandler(store)

		go func() {
			if err := server.RegisterRoutes(); err != nil {
				log.Fatal(err)
			}
		}()

		ticker := time.NewTicker(time.Duration(cfg.CheckInterval) * time.Second)
		defer ticker.Stop()

		for {
			health.Check(cfg, store)
			<-ticker.C
		}
	}
}
